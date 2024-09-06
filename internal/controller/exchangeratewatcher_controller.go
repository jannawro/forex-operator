/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	forexv1alpha1 "github.com/jannawro/forex-operator/api/v1alpha1"
	"github.com/jannawro/forex-operator/internal/forex"
)

// ExchangeRateWatcherReconciler reconciles a ExchangeRateWatcher object
type ExchangeRateWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=forex.jannawro.dev,resources=exchangeratewatchers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=forex.jannawro.dev,resources=exchangeratewatchers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=forex.jannawro.dev,resources=exchangeratewatchers/finalizers,verbs=update
//+kubebuilder:rbac:groups=forex.jannawro.dev,resources=exchangerates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=forex.jannawro.dev,resources=exchangerates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=forex.jannawro.dev,resources=exchangerates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ExchangeRateWatcher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ExchangeRateWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Starting reconciliation loop", "exchangeRateWatcher", prettyPrintNamespacedName(req.NamespacedName))
	watcher := &forexv1alpha1.ExchangeRateWatcher{}
	err := r.Get(ctx, req.NamespacedName, watcher)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			logger.Error(err, "Not found", "exchangeRateWatcher", prettyPrintNamespacedName(req.NamespacedName))
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Unexpected error when fetching", "exchangeRateWatcher", prettyPrintNamespacedName(req.NamespacedName))
		return ctrl.Result{}, err
	}

	client, err := forex.New(watcher.Spec.BaseCurrency, watcher.Spec.TargetCurrencies, 60*time.Second)
	if err != nil {
		logger.Error(err, "Failed creating forex API client")
		return ctrl.Result{}, err
	}
	currencyRates, err := client.GetRates()
	if err != nil {
		logger.Error(err, "Failed fetching exchange rates")
		return ctrl.Result{}, err
	}

	for currency, rate := range currencyRates {
		rateConverted := fmt.Sprintf("%.2f", rate)
		if err := r.createOrUpdateExchangeRate(ctx, watcher, currency, rateConverted); err != nil {
			return ctrl.Result{}, err
		}
	}

	logger.Info("Updating status", "exchangeRateWatcher", prettyPrintNamespacedName(req.NamespacedName))
	watcher.Status.LastChecked = metav1.Now()
	if err := r.Status().Update(ctx, watcher); err != nil {
		logger.Error(err, "Failed updating status", "exchangeRateWatcher", prettyPrintNamespacedName(req.NamespacedName))
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(watcher.Spec.WatchIntervalSeconds) * time.Second}, nil
}

func (r *ExchangeRateWatcherReconciler) createOrUpdateExchangeRate(ctx context.Context, watcher *forexv1alpha1.ExchangeRateWatcher, targetCurrency string, rate string) error {
	logger := log.FromContext(ctx)
	exchangeRate := &forexv1alpha1.ExchangeRate{}

	suffix := base64EncodeAndTrim(watcher.Name)
	exchangeRateName := strings.ToLower(fmt.Sprintf("%s-to-%s-%s", watcher.Spec.BaseCurrency, targetCurrency, suffix))
	namespacedName := types.NamespacedName{
		Namespace: watcher.Namespace,
		Name:      exchangeRateName,
	}

	logger.Info("Starting reconciliation", "exchangeRate", prettyPrintNamespacedName(namespacedName))
	err := r.Get(ctx, namespacedName, exchangeRate)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			logger.Info("Not found", "exchangeRate", prettyPrintNamespacedName(namespacedName))
			newExchangeRate := &forexv1alpha1.ExchangeRate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      namespacedName.Name,
					Namespace: namespacedName.Namespace,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: watcher.APIVersion,
							Kind:       watcher.Kind,
							Name:       watcher.Name,
							UID:        watcher.UID,
							Controller: ptr.To(true),
						},
					},
				},
				Spec: forexv1alpha1.ExchangeRateSpec{
					BaseCurrency:   watcher.Spec.BaseCurrency,
					TargetCurrency: targetCurrency,
				},
				Status: forexv1alpha1.ExchangeRateStatus{
					Rate:        rate,
					LastUpdated: metav1.Now(),
				},
			}
			logger.Info("Creating", "exchangeRate", prettyPrintNamespacedName(namespacedName))
			return r.Create(ctx, newExchangeRate)
		}
		return err
	}

	logger.Info("Updating rate and status", "exchangeRate", prettyPrintNamespacedName(namespacedName))
	exchangeRate.Status.Rate = rate
	exchangeRate.Status.LastUpdated = metav1.Now()
	return r.Status().Update(ctx, exchangeRate)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExchangeRateWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&forexv1alpha1.ExchangeRateWatcher{}).
		Owns(&forexv1alpha1.ExchangeRate{}).
		Complete(r)
}

func base64EncodeAndTrim(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	if len(encoded) > 8 {
		return encoded[:8]
	}
	return encoded
}

func prettyPrintNamespacedName(n types.NamespacedName) string {
	return fmt.Sprint(n.Namespace + "/" + n.Name)
}
