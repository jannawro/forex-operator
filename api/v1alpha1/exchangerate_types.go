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

// ExchangeRateSpec defines the desired state of ExchangeRate
type ExchangeRateSpec struct {
	BaseCurrency   string `json:"baseCurrency,omitempty"`
	TargetCurrency string `json:"targetCurrency,omitempty"`
}

// ExchangeRateStatus defines the observed state of ExchangeRate
type ExchangeRateStatus struct {
	Rate        string      `json:"rate,omitempty"`
	LastUpdated metav1.Time `json:"lastUpdated,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=er
//+kubebuilder:printcolumn:name="Base Currency",type=string,JSONPath=`.spec.baseCurrency`
//+kubebuilder:printcolumn:name="Target Currency",type=string,JSONPath=`.spec.targetCurrency`
//+kubebuilder:printcolumn:name="Rate",type=string,JSONPath=`.status.rate`
//+kubebuilder:printcolumn:name="Last Updated",type=date,JSONPath=`.status.lastUpdated`

// ExchangeRate is the Schema for the exchangerates API
type ExchangeRate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExchangeRateSpec   `json:"spec,omitempty"`
	Status ExchangeRateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExchangeRateList contains a list of ExchangeRate
type ExchangeRateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExchangeRate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExchangeRate{}, &ExchangeRateList{})
}
