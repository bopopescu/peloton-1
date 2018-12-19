package orm

import (
	"code.uber.internal/infra/peloton/storage/objects/base"
)

// ValidObject is a representation of the orm annotations
type ValidObject struct {
	base.Object `cassandra:"name=valid_object, primaryKey=((id), name)"`
	ID          uint64 `column:"name=id"`
	Name        string `column:"name=name"`
	Data        string `column:"name=data"`
}

// InvalidObject1 has primary key as empty
type InvalidObject1 struct {
	base.Object `cassandra:"name=valid_object, primaryKey=()"`
	ID          uint64 `column:"name=id"`
	Name        string `column:"name=name"`
}

// InvalidObject2 has invalid orm tag
type InvalidObject2 struct {
	base.Object `randomstring:"name=valid_object, primaryKey=((id), name)"`
	ID          uint64 `column:"name=id"`
	Name        string `column:"name=name"`
}

// InvalidObject3 has invalid orm tag on ID field
type InvalidObject3 struct {
	base.Object `cassandra:"name=valid_object, primaryKey=((id), name)"`
	ID          uint64 `randomstring:"name=id"`
	Name        string `column:"name=name"`
}

// TestTableFromObject tests creating orm.Table from given base object
// This is meant to test that only entities annotated in a certain format will
// be successfully converted to orm tables
func (suite *ORMTestSuite) TestTableFromObject() {
	_, err := TableFromObject(&ValidObject{})
	suite.NoError(err)

	tt := []base.Object{
		&InvalidObject1{}, &InvalidObject2{}, &InvalidObject3{}}
	for _, t := range tt {
		_, err := TableFromObject(t)
		suite.Error(err)
	}
}

// TestSetObjectFromRow tests setting base object from a row
func (suite *ORMTestSuite) TestSetObjectFromRow() {
	e := &ValidObject{}
	table, err := TableFromObject(e)
	suite.NoError(err)

	table.SetObjectFromRow(e, testRow)
	suite.Equal(e.ID, testRow[0].Value)
	suite.Equal(e.Name, testRow[1].Value)
	suite.Equal(e.Data, testRow[2].Value)
}

// TestGetRowFromObject tests building a row (list of base.Column) from base
// object
func (suite *ORMTestSuite) TestGetRowFromObject() {
	e := &ValidObject{
		ID:   uint64(1),
		Name: "test",
		Data: "testdata",
	}
	table, err := TableFromObject(e)
	suite.NoError(err)

	row := table.GetRowFromObject(e)
	suite.ensureRowsEqual(row, testRow)
}

// TestGetKeyRowFromObject tests getting primary key row (list of primary key
// base.Column) from base object
func (suite *ORMTestSuite) TestGetKeyRowFromObject() {
	e := &ValidObject{
		ID:   uint64(1),
		Name: "test",
		Data: "junk",
	}
	table, err := TableFromObject(e)
	suite.NoError(err)

	keyRow := table.GetKeyRowFromObject(e)
	suite.Equal(e.ID, keyRow[0].Value)
	suite.Equal(e.Name, keyRow[1].Value)
	suite.Equal(len(keyRow), 2)
}