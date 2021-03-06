func TestNewTables(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	{{with .TestSQLDumpGlobPath}}defer dmltest.SQLDumpLoad(t, "{{.}}", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	})(){{end}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	tbls, err := NewTables(ctx, ddl.WithConnPool(db))
	assert.NoError(t, err)

	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{ {{- range $table := .Tables }}"{{ .TableName}}",{{- end}}}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	var ps *pseudo.Service
	ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de",FloatMaxDecimals:6},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		{{- CustomCode "pseudo.MustNewService.Option" -}}
	)

	// TODO run those tests in parallel
	{{- range $table := .Tables }}
	t.Run("{{GoCamel .TableName}}_Entity", func(t *testing.T) {
		ccd := tbls.MustTable(TableName{{GoCamel .TableName}})

		inStmt, err := ccd.Insert().BuildValues().Prepare(ctx) // Do not use Ignore() to suppress DB errors.
		assert.NoError(t, err)
		insArtisan := inStmt.WithArgs()
		defer dmltest.Close(t, inStmt)

		selArtisan := ccd.SelectByPK().WithArgs().ExpandPlaceHolders()

		for i := 0; i < 9; i++ {
			entityIn := new({{GoCamel .TableName}})
			if err := ps.FakeData(entityIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.{{GoCamel .TableName}}_Entity")(insArtisan.Record("", entityIn).ExecContext(ctx))
			insArtisan.Reset()

			entityOut := new({{GoCamel .TableName}})
			rowCount, err := selArtisan.Int64s(lID).Load(ctx, entityOut)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)

			{{- range $col := $table.Columns }}
				{{if $col.IsString -}}
					assert.ExactlyLength(t, {{$col.CharMaxLength.Int64}}, &entityIn.{{$table.GoCamelMaybePrivate $col.Field}}, &entityOut.{{$table.GoCamelMaybePrivate $col.Field}}, "IDX%d: {{$table.GoCamelMaybePrivate $col.Field}} should match", lID)
				{{- else if not $col.IsSystemVersioned -}}
					assert.Exactly(t, entityIn.{{$table.GoCamelMaybePrivate $col.Field}}, entityOut.{{$table.GoCamelMaybePrivate $col.Field}}, "IDX%d: {{$table.GoCamelMaybePrivate $col.Field}} should match", lID)
				{{- end}}
			{{- end}}
		}
	})
	{{- end}}
}
