// Licensed to Typogenerator under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Typogenerator licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hduplooy/gocsv"
	"github.com/namsral/flag"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"golang.org/x/net/idna"

	"Shahriar-Sazid/typogenerator/mapping"
	"Shahriar-Sazid/typogenerator/strategy"
)

var (
	input           = flag.String("s", "mymensingh", "Defines string to alternate")
	permutationOnly = flag.Bool("p", true, "Display permutted domain only")
)

// FuzzResult represents permutations results
type FuzzResult struct {
	StrategyName string   `json:"name" yaml:"name"`
	Domain       string   `json:"domain" yaml:"domain"`
	Permutations []string `json:"permutations" yaml:"permutations"`
}

// Fuzz string using given strategies
func Fuzz(name string, strategies ...strategy.Strategy) ([]FuzzResult, error) {
	return fuzz(name, "", strategies...)
}

// FuzzDomain splits a domain into (TRD + SLD) + TLD and fuzzes using given strategies
func FuzzDomain(domain string, strategies ...strategy.Strategy) ([]FuzzResult, error) {
	parsed, err := publicsuffix.Parse(domain)
	if err != nil {
		return []FuzzResult{}, err
	}

	domain = parsed.SLD
	if parsed.TRD != "" {
		domain = parsed.TRD + "." + domain
	}

	return fuzz(domain, parsed.TLD, strategies...)
}

func fuzz(domain, tld string, strategies ...strategy.Strategy) ([]FuzzResult, error) {
	res := []FuzzResult{}
	var err error

	var domains []string
	for _, s := range strategies {
		if s != nil {
			r := FuzzResult{
				StrategyName: s.GetName(),
				Domain:       domain,
			}

			domains, err = s.Generate(domain, tld)
			if err != nil {
				break
			}

			// Assign permutations to result
			r.Permutations = domains

			// Add result
			res = append(res, r)
		}
	}

	return res, err
}

func init() {
	flag.Parse()
}

func main() {
	all := []strategy.Strategy{
		strategy.DoubleHit(mapping.English),
		strategy.VowelSwap,
		strategy.Similar(mapping.English),
		strategy.Omission,
		strategy.Transposition,
		strategy.Repetition,
		strategy.Replace(mapping.English),
	}

	results, err := Fuzz(*input, all...)
	if err != nil {
		log.Fatal("Unable to generate domains.")
	}

	if !*permutationOnly {
		writer := gocsv.NewWriter(os.Stdout)
		writer.QuoteFields = true
		defer writer.Flush()

		// Write headers
		if err := writer.Write([]string{"strategy", "domain", "permunation", "idna"}); err != nil {
			panic(err)
		}

		for _, r := range results {
			for _, p := range r.Permutations {
				puny, _ := idna.ToASCII(p)
				if err := writer.Write([]string{r.StrategyName, r.Domain, p, puny}); err != nil {
					panic(err)
				}
			}
		}
	} else {
		for _, r := range results {
			for _, p := range r.Permutations {
				fmt.Println(p)
			}
		}
	}
}
