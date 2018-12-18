package loggable

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var db *gorm.DB

type SomeType struct {
	gorm.Model
	Source    string
	FirstName string `gorm-loggable:"true"`
	LastName  string `gorm-loggable:"true"`
	MetaModel
}

type MetaModel struct {
	createdBy string
	LoggableModel
}

func (m MetaModel) Meta() interface{} {
	return struct {
		CreatedBy string
	}{CreatedBy: m.createdBy}
}

//func TestMain(m *testing.M) {
//	database, err := gorm.Open(
//		"postgres",
//		fmt.Sprintf(
//			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
//			"root",
//			"keepitsimple",
//			"localhost",
//			5432,
//			"loggable",
//		),
//	)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//	}
//	database = database.LogMode(true)
//	_, err = Register(database)
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//	}
//	err = database.AutoMigrate(SomeType{}).Error
//	if err != nil {
//		fmt.Println(err)
//		panic(err)
//	}
//	db = database
//	m.Run()
//}

func TestTryModel(t *testing.T) {
	newmodel := SomeType{Source: time.Now().Format(time.Stamp)}
	newmodel.createdBy = "some user"
	err := db.Create(&newmodel).Error
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(newmodel.ID)
	newmodel.Source = "updated field"
	err = db.Model(SomeType{}).Save(&newmodel).Error
	if err != nil {
		t.Fatal(err)
	}
}

func TestIdentityManager(t *testing.T) {
	m := newIdentityManager()

	m1 := SomeType{}
	m2 := &SomeType{}
	m.save(m1, "123")
	m.save(m2, "456")

	v1 := m.get(SomeType{}, "123")
	assert.NotNil(t, v1)
	assert.IsType(t, SomeType{}, v1)

	v2 := m.get(SomeType{}, "456")
	assert.NotNil(t, v2)
	assert.IsType(t, SomeType{}, v2)

	v3 := m.get(SomeType{}, "111")
	assert.Nil(t, v3)
}

func TestIdentifyManagerDiff(t *testing.T) {
	m := newIdentityManager()
	a := SomeType{
		FirstName: "John",
		LastName:  "Doe",
	}

	b := SomeType{
		FirstName: "Jane",
		LastName:  "Doe",
	}

	m.save(a, "111")

	diff := m.diff(b, "111")
	assert.NotNil(t, diff)
}

func TestGetLoggableFieldNames(t *testing.T) {
	m := SomeType{}

	names := getLoggableFieldNames(m)
	assert.NotEmpty(t, names)
	assert.Equal(t, []string{"FirstName", "LastName"}, names)

	pm := SomeType{}
	names = getLoggableFieldNames(pm)
	assert.NotEmpty(t, names)
	assert.Equal(t, []string{"FirstName", "LastName"}, names)
}

func TestComputeDiff(t *testing.T) {
	a := SomeType{
		FirstName: "John",
		LastName:  "Doe",
	}

	b := SomeType{
		FirstName: "Jane",
		LastName:  "Doe",
	}

	diff := computeDiff(a, b)
	assert.Equal(t, diff["FirstName"], "Jane")
	assert.Len(t, diff, 1)

	x := &SomeType{
		FirstName: "John",
		LastName:  "Doe",
	}

	y := SomeType{
		FirstName: "Jane",
		LastName:  "Simpson",
	}

	diff = computeDiff(x, y)
	assert.Equal(t, diff["FirstName"], "Jane")
	assert.Equal(t, diff["LastName"], "Simpson")
}

func TestGenIdentityHash(t *testing.T) {
	a := SomeType{
		FirstName: "John",
		LastName:  "Doe",
	}
	b := SomeType{
		FirstName: "Jane",
		LastName:  "Doe",
	}

	ha := genIdentityHash(a, "123")
	hb := genIdentityHash(b, "123")
	assert.Equal(t, ha, hb)

	x := &SomeType{
		FirstName: "John",
		LastName:  "Doe",
	}
	y := SomeType{
		FirstName: "Jane",
		LastName:  "Doe",
	}
	hx := genIdentityHash(x, "123")
	hy := genIdentityHash(y, "123")
	assert.Equal(t, hx, hy)
}
