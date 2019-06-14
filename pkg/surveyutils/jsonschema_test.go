package surveyutils_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/AlecAivazis/survey.v1/core"

	expect "github.com/Netflix/go-expect"
	"github.com/ghodss/yaml"

	"github.com/jenkins-x/jx/pkg/tests"

	"github.com/stretchr/testify/assert"

	jenkinsv1 "github.com/jenkins-x/jx/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx/pkg/surveyutils"
)

// TODO Figure out how to test selects (affects arrays, enums, validation keywords for arrays)

var timeout = 5 * time.Second

type GeneratedSecret struct {
	Name  string
	Key   string
	Value string
}

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestObjectType(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "objectType.test.schema.json", make(map[string]interface{}), false, false,
			false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for name")
				console.SendLine("cheese")
				console.ExpectEOF()
			})
		assert.Equal(r, `nestedObject:
  anotherNestedObject:
    name: cheese
`, values)
		assert.NoError(r, err)
	})
}

func TestDescriptionAndTitle(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "descriptionAndTitle.test.schema.json", make(map[string]interface{}),
			false,
			false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test explicit question
				console.ExpectString("What is your name?")
				console.SendLine("?")
				// Test explicit description
				console.ExpectString("Enter your name")
				console.SendLine("Pete")
				// Test no description
				console.ExpectString("What is your address?")
				console.SendLine("?")
				// Test no title
				console.ExpectString("Enter a value for country")
				console.SendLine("UK")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `address: '?'
country: UK
name: Pete
`, values)
	})
}

func TestAutoAcceptDefaultValues(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "autoAcceptDefaultValues.test.schema.json", make(map[string]interface{}),
			false, false,
			true, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test explicit question
				//console.ExpectString("What is your name? John Smith [Automatically accepted default value]")
				//console.ExpectEOF()
				// TODO Fix the console test
			})
		assert.Equal(r, `name: John Smith
`, values)
		assert.NoError(r, err)
	})
}

func TestAcceptExisting(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		t.SkipNow()
		// TODO Fix failing test
		values, _, err := GenerateValuesAsYaml(r, "acceptExisting.test.schema.json", map[string]interface{}{
			"name": "John Smith",
		},
			false, false,
			false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test explicit question
				console.ExpectString("What is your name? John Smith [Automatically accepted existing value]")
				console.ExpectEOF()
			})
		assert.Equal(r, `name: John Smith
`, values)
		assert.NoError(r, err)
	})
}

func TestAskExisting(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "askExisting.test.schema.json", map[string]interface{}{
			"name": "John Smith",
		},
			true,
			false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test explicit question
				console.ExpectString("What is your name? [? for help] (John Smith)")
				console.SendLine("")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `name: John Smith
`, values)
	})
}

func TestNoAskAndAutoAcceptDefaultsWithExisting(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		// TODO Fix the flacky console tests and reenable this test again
		t.Skip()
		values, _, err := GenerateValuesAsYaml(r, "noAskAndAutoAcceptDefaultsWithExisting.test.schema.json",
			map[string]interface{}{
				"name":    "John Smith",
				"country": "UK",
			},
			false,
			true, true, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test explicit question
				console.ExpectString("What is your name? John Smith [Automatically accepted existing value]")
				console.ExpectString("Enter a value for country UK [Automatically accepted default value]")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `country: UK
name: John Smith
`, values)
	})
}

func TestIgnoreMissingValues(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "ignoreMissingValues.test.schema.json", make(map[string]interface{}),
			false,
			true,
			false, true,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `{}
`, values)
	})
}

func TestErrorMissingValues(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "ignoreMissingValues.test.schema.json", make(map[string]interface{}),
			false,
			true,
			false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectEOF()
			})
		assert.Error(t, err)
	})
}

func TestDefaultValues(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "defaultValues.test.schema.json", make(map[string]interface{}), false,
			false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test default value
				console.ExpectString("Enter a value for stringValue (UK)")
				console.SendLine("")
				console.ExpectString("Enter a value for booleanValue (y/N)")
				console.SendLine("")
				console.ExpectString("Enter a value for numberValue (123.4)")
				console.SendLine("")
				console.ExpectString("Enter a value for integerValue (123)")
				console.SendLine("")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `booleanValue: false
integerValue: 123
numberValue: 123.4
stringValue: UK
`, values)
	})
}

func TestConstValues(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "constValues.test.schema.json", make(map[string]interface{}), false,
			false,
			false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test default value
				console.ExpectString("Do you want to set stringValue to UK (Y/n)")
				console.SendLine("")
				console.ExpectString("Do you want to set booleanValue to false (Y/n)")
				console.SendLine("")
				console.ExpectString("Do you want to set numberValue to 123.4 (Y/n)")
				console.SendLine("")
				console.ExpectString("Do you want to set integerValue to 123")
				console.SendLine("")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `booleanValue: false
integerValue: 123
numberValue: 123.4
stringValue: UK
`, values)
	})
}

func TestBasicTypesValidation(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "basicTypesValidation.test.schema.json", make(map[string]interface{}), false,
			false,
			false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: unable to convert abc to float64")
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("123.1")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("123.1")
				console.ExpectString("Sorry, your reply was invalid: unable to convert 123.1 to int")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("123")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestBasicTypes(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "basicTypes.test.schema.json", make(map[string]interface{}), false, false,
			false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for booleanValue (y/N)")
				console.SendLine("Y")
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("123.4")
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("hello")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("123")
				console.ExpectEOF()
			})
		assert.Equal(r, `booleanValue: true
integerValue: 123
nullValue: null
numberValue: 123.4
stringValue: hello
`, values)
		assert.NoError(r, err)
	})
}

func TestMultipleOf(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "multipleOf.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("11.1")
				console.ExpectString("Sorry, your reply was invalid: 11.1 cannot be divided by 10")
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("10")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("12")
				console.ExpectString("Sorry, your reply was invalid: 12 cannot be divided by 20")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("20")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestMaximum(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "maximum.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("11.1")
				console.ExpectString("Sorry, your reply was invalid: 11.1 is not less than or equal to 10.1")
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("1")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("21")
				console.ExpectString("Sorry, your reply was invalid: 21 is not less than or equal to 20")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("2")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestExclusiveMaximum(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "exclusiveMaximum.test.schema.json", make(map[string]interface{}), false, false, false,
			false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("10.1")
				console.ExpectString("Sorry, your reply was invalid: 10.1 is not less than 10.1")
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("1")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("20")
				console.ExpectString("Sorry, your reply was invalid: 20 is not less than 20")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("2")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestMinimum(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "minimum.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("9.1")
				console.ExpectString("Sorry, your reply was invalid: 9.1 is not greater than or equal to 10.1")
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("11")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("19")
				console.ExpectString("Sorry, your reply was invalid: 19 is not greater than or equal to 20")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("21")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestExclusiveMinimum(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "exclusiveMinimum.test.schema.json", make(map[string]interface{}), false, false, false,
			false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("10.1")
				console.ExpectString("Sorry, your reply was invalid: 10.1 is not greater than 10.1")
				console.ExpectString("Enter a value for numberValue")
				console.SendLine("11")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("20")
				console.ExpectString("Sorry, your reply was invalid: 20 is not greater than 20")
				console.ExpectString("Enter a value for integerValue")
				console.SendLine("21")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestMaxLength(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "maxLength.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("iamlongerthan10")
				console.ExpectString("Sorry, your reply was invalid: value is too long. Max length is 10")
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("short")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestMinLength(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "minLength.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("short")
				console.ExpectString("Sorry, your reply was invalid: value is too short. Min length is 10")
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("iamlongerthan10")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestPattern(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "pattern.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("HELLO")
				console.ExpectString("Sorry, your reply was invalid: HELLO does not match [0-9]")
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("123")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestRequired(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "required.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("")
				console.ExpectString("Sorry, your reply was invalid: Value is required")
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("Hello")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestIfThen(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "ifThenElse.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for enablePersistentStorage")
				console.SendLine("Y")
				console.ExpectString("Enter a value for databaseConnectionUrl")
				console.SendLine("abc")
				console.ExpectString("Enter a value for databaseUsername")
				console.SendLine("wensleydale")
				console.ExpectString("Enter a value for databasePassword")
				console.SendLine("cranberries")
				console.ExpectString(" ***********")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `databaseConnectionUrl: abc
databasePassword:
  kind: Secret
  name: databasepassword-secret
databaseUsername: wensleydale
enablePersistentStorage: true
`, values)
	})
}

func TestIfElse(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "ifThenElse.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for enablePersistentStorage")
				console.SendLine("N")
				console.ExpectString("Enter a value for enableInMemoryDB")
				console.SendLine("N")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `enableInMemoryDB: false
enablePersistentStorage: false
`, values)
	})
}

func TestIfElseNested(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "ifThenElseNested.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for enablePersistentStorage")
				console.SendLine("N")
				console.ExpectString("Enter a value for enableInMemoryDB")
				console.SendLine("Y")
				console.ExpectString("Enter a value for nestedString")
				console.SendLine("Test")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `nestedObject:
  enableInMemoryDB: true
  enablePersistentStorage: false
  nestedString: Test
`, values)
	})
}

func TestIfElseWithDefaults(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "ifThenElse.test.schema.json", make(map[string]interface{}), false, false, true, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for enablePersistentStorage")
				console.SendLine("N")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `enableInMemoryDB: true
enablePersistentStorage: false
`, values)
	})
}

func TestAllOf(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "AllOfIf.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for enablePersistentStorage")
				console.SendLine("Y")
				console.ExpectString("Enter a value for databaseConnectionUrl")
				console.SendLine("abc")
				console.ExpectString("Enter a value for databaseUsername")
				console.SendLine("wensleydale")
				console.ExpectString("Enter a value for databasePassword")
				console.SendLine("cranberries")
				console.ExpectString(" ***********")
				console.ExpectString("Enter a value for enableCheese")
				console.SendLine("Y")
				console.ExpectString("Enter a value for cheeseType")
				console.SendLine("Stilton")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `cheeseType: Stilton
databaseConnectionUrl: abc
databasePassword:
  kind: Secret
  name: databasepassword-secret
databaseUsername: wensleydale
enableCheese: true
enablePersistentStorage: true
`, values)
	})
}

func TestAllOfThen(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, _, err := GenerateValuesAsYaml(r, "AllOfIf.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				console.ExpectString("Enter a value for enablePersistentStorage")
				console.SendLine("Y")
				console.ExpectString("Enter a value for databaseConnectionUrl")
				console.SendLine("abc")
				console.ExpectString("Enter a value for databaseUsername")
				console.SendLine("wensleydale")
				console.ExpectString("Enter a value for databasePassword")
				console.SendLine("cranberries")
				console.ExpectString(" ***********")
				console.ExpectString("Enter a value for enableCheese")
				console.SendLine("N")
				console.ExpectString("Enter a value for iDontLikeCheese")
				console.SendLine("Y")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
		assert.Equal(r, `databaseConnectionUrl: abc
databasePassword:
  kind: Secret
  name: databasepassword-secret
databaseUsername: wensleydale
enableCheese: false
enablePersistentStorage: true
iDontLikeCheese: true
`, values)
	})
}

func TestMinProperties(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "minProperties.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("")
				console.ExpectString("Enter a value for stringValue1")
				console.SendLine("")
				console.ExpectString("Sorry, your reply was invalid: nestedObject has less than 1 items, has []")
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("abc")
				console.ExpectString("Enter a value for stringValue1")
				console.SendLine("def")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestMaxProperties(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "maxProperties.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("abc")
				console.ExpectString("Enter a value for stringValue1")
				console.SendLine("def")
				console.ExpectString("Sorry, your reply was invalid: nestedObject has more than 1 items, " +
					"has [stringValue stringValue1]")
				console.ExpectString("Enter a value for stringValue")
				console.SendLine("abc")
				console.ExpectString("Enter a value for stringValue1")
				console.SendLine("")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestDateTime(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "dateTime.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for dateTimeValue")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: abc is not a RFC 3339 date-time formatted string, " +
					"it should be like 2006-01-02T15:04:05Z07:00")
				console.ExpectString("Enter a value for dateTimeValue")
				console.SendLine("2006-01-02T15:04:05-07:00")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestDate(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "date.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for dateValue")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: abc is not a RFC 3339 full-date formatted string, " +
					"it should be like 2006-01-02")
				console.ExpectString("Enter a value for dateValue")
				console.SendLine("2006-01-02")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestTime(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "time.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for timeValue")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: abc is not a RFC 3339 full-time formatted string, " +
					"it should be like 15:04:05Z07:00")
				console.ExpectString("Enter a value for timeValue")
				console.SendLine("15:04:05-07:00")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestPassword(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, secrets, err := GenerateValuesAsYaml(r, "password.test.schema.json", make(map[string]interface{}), false,
			false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for passwordValue")
				console.SendLine("abc")
				console.ExpectEOF()
			})
		assert.Equal(r, `passwordValue:
  kind: Secret
  name: passwordvalue-secret
`, values)
		assert.Contains(r, secrets, &GeneratedSecret{
			Name:  "passwordvalue-secret",
			Value: "abc",
			Key:   "password",
		})
		assert.NoError(r, err)
	})
}

func TestToken(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		values, secrets, err := GenerateValuesAsYaml(r, "token.test.schema.json", make(map[string]interface{}), false,
			false,
			false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for tokenValue")
				console.SendLine("abc")
				console.ExpectEOF()
			})
		assert.Equal(r, `tokenValue:
  kind: Secret
  name: tokenvalue-secret
`, values)
		assert.Contains(r, secrets, &GeneratedSecret{
			Name:  "tokenvalue-secret",
			Value: "abc",
			Key:   "token",
		})
		assert.NoError(r, err)
	})
}

func TestEmail(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "email.test.schema.json", make(map[string]interface{}), false, false, false,
			false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for emailValue")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: abc is not a RFC 5322 address, " +
					"it should be like Barry Gibb <bg@example.com>")
				console.ExpectString("Enter a value for emailValue")
				console.SendLine("Maurice Gibb <mg@example.com>")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestIdnEmail(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "idnemail.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for emailValue")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: abc is not a RFC 5322 address, " +
					"it should be like Barry Gibb <bg@example.com>")
				console.ExpectString("Enter a value for emailValue")
				console.SendLine("Maurice Gibb <mg@example.com>")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestHostname(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "hostname.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for hostnameValue")
				console.SendLine("*****")
				console.ExpectString("Sorry, your reply was invalid: ***** is not a RFC 1034 hostname, " +
					"it should be like example.com")
				console.ExpectString("Enter a value for hostnameValue")
				console.SendLine("example.com")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestIdnHostname(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "idnhostname.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for hostnameValue")
				console.SendLine("*****")
				console.ExpectString("Sorry, your reply was invalid: ***** is not a RFC 1034 hostname, " +
					"it should be like example.com")
				console.ExpectString("Enter a value for hostnameValue")
				console.SendLine("example.com")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestIpv4(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "ipv4.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for ipv4Value")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: abc is not a RFC 2673 IPv4 Address, " +
					"it should be like 127.0.0.1")
				console.ExpectString("Enter a value for ipv4Value")
				console.SendLine("127.0.0.1")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestIpv6(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "ipv6.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for ipv6Value")
				console.SendLine("abc")
				console.ExpectString("Sorry, your reply was invalid: abc is not a RFC 4291 IPv6 address, " +
					"it should be like ::1")
				console.ExpectString("Enter a value for ipv6Value")
				console.SendLine("::1")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestUri(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "uri.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for uriValue")
				console.SendLine("*****")
				console.ExpectString("Sorry, your reply was invalid: ***** is not a RFC 3986 URI")
				console.ExpectString("Enter a value for uriValue")
				console.SendLine("https://example.com")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestUriReference(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "uriReference.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for uriReferenceValue")
				console.SendLine("http$$://foo")
				console.ExpectString("Sorry, your reply was invalid: http$$://foo is not a RFC 3986 URI reference")
				console.ExpectString("Enter a value for uriReferenceValue")
				console.SendLine("../resource.txt")
				console.ExpectEOF()
			})
		assert.NoError(r, err)
	})
}

func TestJSONPointer(t *testing.T) {
	tests.SkipForWindows(t, "go-expect does not work on windows")
	tests.Retry(t, 5, time.Second*10, func(r *tests.R) {
		_, _, err := GenerateValuesAsYaml(r, "jsonPointer.test.schema.json", make(map[string]interface{}), false, false, false, false,
			func(console *tests.ConsoleWrapper, donec chan struct{}) {
				defer close(donec)
				// Test boolean type
				console.ExpectString("Enter a value for jsonPointerValue")
				console.SendLine("~")
				console.ExpectString("Sorry, your reply was invalid: ~ is not a RFC 6901 JSON pointer")
				console.ExpectString("Enter a value for jsonPointerValue")
				console.SendLine("/abc")
				console.ExpectEOF()
			})
		assert.NoError(r, err)

	})
}

func GenerateValuesAsYaml(r *tests.R, schemaName string, existingValues map[string]interface{},
	askExisting bool, noAsk bool, autoAcceptDefaults bool, ignoreMissingValues bool, answerQuestions func(
		console *tests.
			ConsoleWrapper, donec chan struct{})) (string, []*GeneratedSecret, error) {

	//t.Parallel()
	secrets := make([]*GeneratedSecret, 0)
	console := tests.NewTerminal(r, &timeout)
	defer console.Cleanup()
	options := surveyutils.JSONSchemaOptions{
		Out:                 console.Out,
		In:                  console.In,
		OutErr:              console.Err,
		AskExisting:         askExisting,
		AutoAcceptDefaults:  autoAcceptDefaults,
		NoAsk:               noAsk,
		IgnoreMissingValues: ignoreMissingValues,

		CreateSecret: func(name string, key string, value string, passthrough bool) (interface{}, error) {
			secrets = append(secrets, &GeneratedSecret{
				Name:  name,
				Value: value,
				Key:   key,
			})
			return &jenkinsv1.ResourceReference{
				Name: name,
				Kind: "Secret",
			}, nil
		},
	}
	data, err := ioutil.ReadFile(filepath.Join("test_data", schemaName))
	assert.NoError(r, err)

	// Test interactive IO
	donec := make(chan struct{})
	go answerQuestions(console, donec)
	assert.NoError(r, err)
	result, runErr := options.GenerateValues(
		data,
		existingValues)
	console.Close()
	<-donec
	yaml, err := yaml.JSONToYAML(result)
	r.Logf(expect.StripTrailingEmptyLines(console.CurrentState()))
	assert.NoError(r, err)
	return string(yaml), secrets, runErr
}
