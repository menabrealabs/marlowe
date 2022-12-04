package language_test

import (
	"testing"

	assert "github.com/menabrealabs/marlowe/assertion"
	c "github.com/menabrealabs/marlowe/v1/language/core"
	ext "github.com/menabrealabs/marlowe/v1/language/extended"
)

func TestTypes_WhenContract(t *testing.T) {
	// Should generate JSON:
	// {"when":[{"then":"close","case":{"for_choice":{"choice_owner":{"role_token":"creditor"},"choice_name":"option"},"choose_between":[{"to":2,"from":1}]}}],"timeout_continuation":"close","timeout":1668250824063}

	contract := c.When{
		Cases: []c.Case{
			{
				Action: c.Choice{
					ChoiceId: c.ChoiceId{
						Name:  "option",
						Owner: c.Role{"creditor"},
					},
					Bounds: []c.Bound{
						{
							Upper: 3,
							Lower: 2,
						},
					},
				},
				Then: c.Close,
			},
		},
		Timeout: ext.TimeParam("deadline"),
		Then:    c.Close,
	}

	assert.Json(t, contract, `{"when":[{"case":{"for_choice":{"choice_name":"option","choice_owner":{"role_token":"creditor"}},"choose_between":[{"from":3,"to":2}]},"then":"close"}],"timeout":1668250824063,"timeout_continuation":"close"}`)
}
