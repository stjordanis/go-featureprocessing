package examplemodule

import (
	"encoding/json"
	"testing"

	. "github.com/nikolaydubina/go-featureprocessing/transformers"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeFeatureTransformerReadme(t *testing.T) {
	t.Run("transform", func(t *testing.T) {
		employee := Employee{
			Age:         22,
			Salary:      1000.0,
			Kids:        2,
			Weight:      85.1,
			Height:      160.0,
			City:        "Pangyo",
			Car:         "Tesla",
			Income:      9000.1,
			SecretValue: 42,
			Description: "large text fields are not a problem neither, tf-idf can help here too! more advanced NLP will be added later!",
		}

		tr := EmployeeFeatureTransformer{
			Salary: MinMaxScaler{Min: 500, Max: 900},
			Kids:   MaxAbsScaler{Max: 4},
			Weight: StandardScaler{Mean: 60, STD: 25},
			Height: QuantileScaler{Quantiles: []float64{20, 100, 110, 120, 150}, NQuantiles: 5},
			City:   OneHotEncoder{Values: []string{"Pangyo", "Seoul", "Daejeon", "Busan"}},
			Car:    OrdinalEncoder{Mapping: map[string]float64{"Tesla": 1, "BMW": 90000}},
			Income: KBinsDiscretizer{QuantileScaler: QuantileScaler{Quantiles: []float64{1000, 1100, 2000, 3000, 10000}, NQuantiles: 5}},
			Description: TfIdfVectorizer{
				NumDocuments:    2,
				DocCount:        map[int]int{0: 1, 1: 2, 2: 2},
				CountVectorizer: CountVectorizer{Mapping: map[string]int{"text": 0, "problem": 1, "help": 2}, Separator: " "},
			},
		}

		features := tr.Transform(&employee)
		expected := []float64{22, 1, 0.5, 1.0039999999999998, 1, 1, 0, 0, 0, 1, 5, 0.7674945674619879, 0.4532946552278861, 0.4532946552278861}
		assert.Equal(t, expected, features)
	})

	t.Run("fit", func(t *testing.T) {
		employee := []Employee{
			{
				Age:         22,
				Salary:      500.0,
				Kids:        2,
				Weight:      50,
				Height:      160.0,
				City:        "Pangyo",
				Car:         "Tesla",
				Income:      9000.1,
				SecretValue: 42,
				Description: "text problem help",
			},
			{
				Age:         10,
				Salary:      900.0,
				Kids:        0,
				Weight:      10,
				Height:      120.0,
				City:        "Seoul",
				Car:         "BMW",
				Income:      420.1,
				Description: "problem help",
			},
		}

		tr := EmployeeFeatureTransformer{}
		tr.Height.NQuantiles = 5
		tr.Income.NQuantiles = 5
		tr.Fit(employee)

		trExpected := EmployeeFeatureTransformer{
			Salary: MinMaxScaler{Min: 500, Max: 900},
			Kids:   MaxAbsScaler{Max: 2},
			Weight: StandardScaler{Mean: 30, STD: 28.284271247461902},
			Height: QuantileScaler{Quantiles: []float64{120, 160}, NQuantiles: 2},
			City:   OneHotEncoder{Values: []string{"Pangyo", "Seoul"}},
			Car:    OrdinalEncoder{Mapping: map[string]float64{"Tesla": 1, "BMW": 2}},
			Income: KBinsDiscretizer{QuantileScaler: QuantileScaler{Quantiles: []float64{420.1, 9000.1}, NQuantiles: 2}},
			Description: TfIdfVectorizer{
				NumDocuments:    2,
				DocCount:        map[int]int{0: 1, 1: 2, 2: 2},
				CountVectorizer: CountVectorizer{Mapping: map[string]int{"text": 0, "problem": 1, "help": 2}, Separator: " "},
			},
		}

		assert.Equal(t, trExpected, tr)
	})

	t.Run("serialize transformer", func(t *testing.T) {
		tr := EmployeeFeatureTransformer{
			Salary: MinMaxScaler{Min: 500, Max: 900},
			Kids:   MaxAbsScaler{Max: 4},
			Weight: StandardScaler{Mean: 60, STD: 25},
			Height: QuantileScaler{Quantiles: []float64{20, 100, 110, 120, 150}, NQuantiles: 5},
			City:   OneHotEncoder{Values: []string{"Pangyo", "Seoul", "Daejeon", "Busan"}},
			Car:    OrdinalEncoder{Mapping: map[string]float64{"Tesla": 1, "BMW": 90000}},
			Income: KBinsDiscretizer{QuantileScaler: QuantileScaler{Quantiles: []float64{1000, 1100, 2000, 3000, 10000}, NQuantiles: 5}},
			Description: TfIdfVectorizer{
				NumDocuments:    2,
				DocCount:        map[int]int{0: 1, 1: 2, 2: 2},
				CountVectorizer: CountVectorizer{Mapping: map[string]int{"text": 0, "problem": 1, "help": 2}, Separator: " "},
			},
		}

		output, err := json.MarshalIndent(tr, "", "    ")
		outputStr := string(output)
		expected := `{
    "Age": {},
    "Salary": {
        "Min": 500,
        "Max": 900
    },
    "Kids": {
        "Max": 4
    },
    "Weight": {
        "Mean": 60,
        "STD": 25
    },
    "Height": {
        "Quantiles": [
            20,
            100,
            110,
            120,
            150
        ],
        "NQuantiles": 5
    },
    "City": {
        "Values": [
            "Pangyo",
            "Seoul",
            "Daejeon",
            "Busan"
        ]
    },
    "Car": {
        "Mapping": {
            "BMW": 90000,
            "Tesla": 1
        }
    },
    "Income": {
        "Quantiles": [
            1000,
            1100,
            2000,
            3000,
            10000
        ],
        "NQuantiles": 5
    },
    "Description": {
        "Mapping": {
            "help": 2,
            "problem": 1,
            "text": 0
        },
        "Separator": " ",
        "DocCount": {
            "0": 1,
            "1": 2,
            "2": 2
        },
        "NumDocuments": 2,
        "Normalizer": {}
    }
}`
		assert.Nil(t, err)
		assert.Equal(t, expected, outputStr)
	})
}
