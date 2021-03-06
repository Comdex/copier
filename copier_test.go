package copier

import (
	"encoding/json"

	"reflect"

	"testing"
)

type User struct {
	Name  string
	Role  string
	Age   int32
	Notes []string
	flags []byte
}

func (user *User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Age       int32
	EmployeID int64
	DoubleAge int32
	SuperRule string
	Notes     []string
	flags     []byte
}

type Base struct {
	BaseField1 int
	BaseField2 int
}

type HaveEmbed struct {
	EmbedField1 int
	EmbedField2 int
	Base
}

type TypeStruct1 struct {
	Field1 	string
	Field2  string
}

type TypeStruct2 struct {
	Field1 	int
	Field2  string
}

type TypeStruct3 struct {
	Field1 	interface{}
	Field2  string
}

type TypeStruct4 struct {
	field1 	int
	Field2  string
}

type TypeStruct5 struct {
	field1 	string
	Field2  string
}

func (t *TypeStruct4) Field1(i int) {
	t.field1 = i
}

func (t *TypeStruct5) Field1(i interface{}) {
	if v, ok := i.(string); ok {
		t.field1 = v
	}
}

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

func TestEmbedded(t *testing.T) {
	base := Base{}
	embeded := HaveEmbed{}
	embeded.BaseField1 = 1
	embeded.BaseField2 = 2
	embeded.EmbedField1 = 3
	embeded.EmbedField2 = 4

	Copy(&base, &embeded)
	
	if base.BaseField1 != 1 {
		t.Error("Embedded fields not copied")
	}
}

func TestDifferentType(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("The copy did panic")
        }
    }()

	ts := &TypeStruct1 {
		Field1: "str1",
		Field2: "str2",
	}

	ts2 := &TypeStruct2 {}

	Copy(ts2, ts)
}

func TestDifferentTypeMethod(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("The copy did panic")
        }
    }()

	ts := &TypeStruct1 {
		Field1: "str1",
		Field2: "str2",
	}

	ts4 := &TypeStruct4 {}

	Copy(ts4, ts)
}

func TestAssignableType(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("The copy did panic")
        }
    }()

	ts := &TypeStruct1 {
		Field1: "str1",
		Field2: "str2",
	}

	ts3 := &TypeStruct3 {}

	Copy(ts3, ts)

	if v, ok := ts3.Field1.(string); !ok {
		t.Error("Assign to interface{} type did not succeed")
	} else if v != "str1" {
		t.Error("String haven't been copied correctly")
	}
}

func TestAssignableTypeMethod(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("The copy did panic")
        }
    }()

	ts := &TypeStruct1 {
		Field1: "str1",
		Field2: "str2",
	}

	ts5 := &TypeStruct5 {}

	Copy(ts5, ts)


	if ts5.field1 != "str1" {
		t.Error("String haven't been copied correctly through method")
	}
}

func TestCopyStruct(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}, flags: []byte{'x'}}
	employee := Employee{}

	Copy(&employee, &user)

	if employee.Name != "Jinzhu" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employee.Age != 18 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employee.DoubleAge != 36 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employee.SuperRule != "Super Admin" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	if !reflect.DeepEqual(employee.Notes, []string{"hello world"}) {
		t.Errorf("Copy a map")
	}

	user.Notes = append(user.Notes, "welcome")
	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("User's Note should be changed")
	}

	if !reflect.DeepEqual(employee.Notes, []string{"hello world"}) {
		t.Errorf("Employee's Note should not be changed")
	}

	employee.Notes = append(employee.Notes, "golang")
	if !reflect.DeepEqual(employee.Notes, []string{"hello world", "golang"}) {
		t.Errorf("Employee's Note should be changed")
	}

	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("Employee's Note should not be changed")
	}
}

func TestCopySlice(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	users := []User{{Name: "jinzhu 2", Age: 30, Role: "Dev"}}
	employees := []Employee{}

	Copy(&employees, &user)
	if len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	}

	Copy(&employees, &users)
	if len(employees) != 2 {
		t.Errorf("Should have two elems when copy additional slice to slice")
	}

	if employees[0].Name != "Jinzhu" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employees[0].Age != 18 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[0].DoubleAge != 36 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employees[0].SuperRule != "Super Admin" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	if employees[1].Name != "jinzhu 2" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employees[1].Age != 30 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[1].DoubleAge != 60 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employees[1].SuperRule != "Super Dev" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	employee := employees[0]
	user.Notes = append(user.Notes, "welcome")
	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("User's Note should be changed")
	}

	if !reflect.DeepEqual(employee.Notes, []string{"hello world"}) {
		t.Errorf("Employee's Note should not be changed")
	}

	employee.Notes = append(employee.Notes, "golang")
	if !reflect.DeepEqual(employee.Notes, []string{"hello world", "golang"}) {
		t.Errorf("Employee's Note should be changed")
	}

	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("Employee's Note should not be changed")
	}
}

func TestCopySliceWithPtr(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	user2 := &User{Name: "jinzhu 2", Age: 30, Role: "Dev"}
	users := []*User{user2}
	employees := []*Employee{}

	Copy(&employees, &user)
	if len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	}

	Copy(&employees, &users)
	if len(employees) != 2 {
		t.Errorf("Should have two elems when copy additional slice to slice")
	}

	if employees[0].Name != "Jinzhu" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employees[0].Age != 18 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[0].DoubleAge != 36 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employees[0].SuperRule != "Super Admin" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	if employees[1].Name != "jinzhu 2" {
		t.Errorf("Name haven't been copied correctly.")
	}
	if employees[1].Age != 30 {
		t.Errorf("Age haven't been copied correctly.")
	}
	if employees[1].DoubleAge != 60 {
		t.Errorf("Copy copy from method doesn't work")
	}
	if employees[1].SuperRule != "Super Dev" {
		t.Errorf("Copy Attributes should support copy to method")
	}

	employee := employees[0]
	user.Notes = append(user.Notes, "welcome")
	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("User's Note should be changed")
	}

	if !reflect.DeepEqual(employee.Notes, []string{"hello world"}) {
		t.Errorf("Employee's Note should not be changed")
	}

	employee.Notes = append(employee.Notes, "golang")
	if !reflect.DeepEqual(employee.Notes, []string{"hello world", "golang"}) {
		t.Errorf("Employee's Note should be changed")
	}

	if !reflect.DeepEqual(user.Notes, []string{"hello world", "welcome"}) {
		t.Errorf("Employee's Note should not be changed")
	}
}

func BenchmarkCopyStruct(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		Copy(&Employee{}, &user)
	}
}

func BenchmarkNamaCopy(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		employee := &Employee{
			Name:      user.Name,
			Age:       user.Age,
			DoubleAge: user.DoubleAge(),
			Notes:     user.Notes,
		}
		employee.Role(user.Role)
	}
}

func BenchmarkJsonMarshalCopy(b *testing.B) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	for x := 0; x < b.N; x++ {
		data, _ := json.Marshal(user)
		var employee Employee
		json.Unmarshal(data, &employee)
		employee.DoubleAge = user.DoubleAge()
		employee.Role(user.Role)
	}
}
