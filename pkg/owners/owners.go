/*
Copyright © 2023 maxgio92 me@maxgio.me

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

package owners

// Handle represents a Github user handle.
type Handle string

// Owners represents a Github OWNERS flie.
type Owners struct {

	// Approvers is a list of users with ability to approve pull requests.
	Approvers []string `json:"approvers"`

	// EmeritusApprovers is a list of users which served for time as approvers.
	EmeritusApprovers []string `json:"emeritus_approvers"`

	// Reviewers is a list of users with ability to review pull requests.
	Reviewers []string `json:"reviewers"`
}

// New returns a new ONWERS content structure.
func New() *Owners {
	return &Owners{}
}
