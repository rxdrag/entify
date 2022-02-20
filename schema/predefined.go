package schema

var MetaEntity = EntityMeta{
	Uuid:       "META_ENTITY_UUID",
	Name:       "Meta",
	TableName:  "meta",
	EntityType: Entity_NORMAL,
	Columns: []ColumnMeta{
		{
			Uuid: "META_COLUMN_ID_UUID",
			Type: COLUMN_ID,
			Name: "id",
		},
		{
			Uuid: "META_COLUMN_VERSION_UUID",
			Type: COLUMN_STRING,
			Name: "version",
		},
		{
			Uuid: "META_COLUMN_CONTENT_UUID",
			Type: COLUMN_STRING,
			Name: "content",
		},
		{
			Uuid: "META_COLUMN_INT_TEST_UUID",
			Type: COLUMN_INT,
			Name: "int_test",
		},
		{
			Uuid: "META_COLUMN_FLOAT_TEST_UUID",
			Type: COLUMN_FLOAT,
			Name: "float_test",
		},
		// {
		// 	Uuid: "_META_COLUMN_PUBLISHED_UUID",
		// 	Type: COLUMN_BOOLEAN,
		// 	Name: "published",
		// },
		{
			Uuid: "META_COLUMN_PUBLISHED_AT_UUID",
			Type: COLUMN_DATE,
			Name: "publishedAt",
		},
		{
			Uuid: "META_COLUMN_CREATED_AT_UUID",
			Type: COLUMN_DATE,
			Name: "createdAt",
		},
		{
			Uuid: "META_COLUMN_UPDATED_AT_UUID",
			Type: COLUMN_DATE,
			Name: "updatedAt",
		},
	},
}
