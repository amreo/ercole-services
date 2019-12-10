/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
package cmd

func init() {
	searchPatchAdvisorsCmd := simpleAPIRequestCommand("search-patch-advisors",
		"Search current patch advisors",
		`search-patch-advisors search the most matching patch advisors to the arguments`,
		true, false, true, true, true, true, false,
		"/patch-advisors",
		"Failed to search patch advisors data: %v\n",
		"Failed to search patch advisors data(Status: %d): %s\n",
	)

	apiCmd.AddCommand(searchPatchAdvisorsCmd)
}