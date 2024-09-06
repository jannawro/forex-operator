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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make" to regenerate code after modifying this file

// ExchangeRateWatcherSpec defines the desired state of ExchangeRateWatcher
type ExchangeRateWatcherSpec struct {
	BaseCurrency         string   `json:"baseCurrency,omitempty"`
	TargetCurrencies     []string `json:"targetCurrencies,omitempty"`
	WatchIntervalSeconds int      `json:"watchIntervalSeconds,omitempty"`
}

// ExchangeRateWatcherStatus defines the observed state of ExchangeRateWatcher
type ExchangeRateWatcherStatus struct {
	LastChecked metav1.Time `json:"lastChecked,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=erw

// ExchangeRateWatcher is the Schema for the exchangeratewatchers API
type ExchangeRateWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExchangeRateWatcherSpec   `json:"spec,omitempty"`
	Status ExchangeRateWatcherStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExchangeRateWatcherList contains a list of ExchangeRateWatcher
type ExchangeRateWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExchangeRateWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExchangeRateWatcher{}, &ExchangeRateWatcherList{})
}
