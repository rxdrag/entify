package meta

import (
	"encoding/json"
	"testing"
)

func TestModifyEntityName(t *testing.T) {
	oldData := `
	{
		"entities": [
			{
				"name": "User",
				"uuid": "3e9ae743-de18-4b0c-a77e-3726be4049a8",
				"columns": [
					{
						"name": "id",
						"type": "ID",
						"uuid": "6758ae89-1e2c-462d-907c-a54baf6cf6fd",
						"primary": true
					},
					{
						"name": "newColumn1",
						"type": "String",
						"uuid": "e1afb0c4-5eee-40f3-8c34-3ce15746877b"
					}
				],
				"entityType": "Normal"
			}
		],
		"relations": []
	}
	`

	newData := `
	{
		"entities": [
			{
				"name": "User2",
				"uuid": "3e9ae743-de18-4b0c-a77e-3726be4049a8",
				"columns": [
					{
						"name": "id",
						"type": "ID",
						"uuid": "6758ae89-1e2c-462d-907c-a54baf6cf6fd",
						"primary": true
					},
					{
						"name": "newColumn1",
						"type": "String",
						"uuid": "e1afb0c4-5eee-40f3-8c34-3ce15746877b"
					}
				],
				"entityType": "Normal"
			}
		],
		"relations": []
	}
	`

	oldM := MetaContent{}
	json.Unmarshal([]byte(oldData), &oldM)
	newM := MetaContent{}
	json.Unmarshal([]byte(newData), &newM)
	newM.Validate()
	diff := CreateDiff(&oldM, &newM)

	if len(diff.ModifiedTables) != 1 {
		t.Errorf("Diffent table number is %d ,not 1", len(diff.ModifiedTables))
	}

	if diff.oldContent.Tables()[0].Name != "user" {
		t.Errorf("Old name is %s, not expected user", diff.oldContent.Tables()[0].Name)
	}

	if diff.newContent.Tables()[0].Name != "user2" {
		t.Errorf("Old name is %s, not expected user2", diff.newContent.Tables()[0].Name)
	}
}

func TestModifiedTableName(t *testing.T) {
	diff := CreateDiff(
		&MetaContent{
			Entities: []Entity{
				{
					Name: "OldName",
				},
			},
		},
		&MetaContent{
			Entities: []Entity{
				{
					Name: "NewName",
				},
			},
		},
	)

	if len(diff.ModifiedTables) != 1 {
		t.Error("Cereate entity name modify diff error, diff number error")
	}

	if diff.ModifiedTables[0].OldTable.Name != "old_name" {
		t.Error("Cereate entity name modify diff error, old name error")
	}

	if diff.ModifiedTables[0].NewTable.Name != "new_name" {
		t.Error("Cereate entity name modify diff error, new name error")
	}

}