//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

package config

func checkRequest(_ *Request) {

}

// SanityCheck verifies a scenario configuration.
func SanityCheck(scenarioFile string) error {

	c := New()

	sc, err := c.ParseFile(scenarioFile)
	if err != nil {
		return err
	}

	for i := range sc.Sequence.Requests {
		checkRequest(&sc.Sequence.Requests[i])
	}

	return nil
}
