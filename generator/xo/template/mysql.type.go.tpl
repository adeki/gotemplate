{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}
type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Name }} {{ retype .Type }} `json:"{{ .Col.ColumnName }}" db:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
{{- end }}
}

{{ if .PrimaryKey }}

// Insert inserts the {{ .Name }} to the database.
func ({{ $short }} *{{ .Name }}) Insert(db XODB) error {
	var err error

{{ if hasfield .Fields "CreatedAt" }}
    {{ $short }}.CreatedAt = time.Now()
{{ end }}
{{ if hasfield .Fields "UpdatedAt" }}
    {{ $short }}.UpdatedAt = time.Now()
{{ end }}


{{ if .Table.ManualPk  }}
	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO {{ $table }} (` +
		`{{ colnames .Fields }}` +
		`) VALUES (` +
		`{{ colvals .Fields }}` +
		`)`

	// run query
	XOLog(sqlstr, {{ fieldnames .Fields $short }})
	_, err = db.Exec(sqlstr, {{ fieldnames .Fields $short }})
	if err != nil {
		return err
	}

{{ else }}
	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO {{ $table }} (` +
		`{{ colnames .Fields .PrimaryKey.Name }}` +
		`) VALUES (` +
		`{{ colvals .Fields .PrimaryKey.Name }}` +
		`)`

	// run query
	XOLog(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	res, err := db.Exec(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }})
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	{{ $short }}.{{ .PrimaryKey.Name }} = {{ .PrimaryKey.Type }}(id)
{{ end }}

	return nil
}

{{ if ne (fieldnamesmulti .Fields $short .PrimaryKeyFields) "" }}
	// Update updates the {{ .Name }} in the database.
	func ({{ $short }} *{{ .Name }}) Update(db XODB) error {
		var err error

        {{ if hasfield .Fields "UpdatedAt" }}
            {{ $short }}.UpdatedAt = time.Now()
        {{ end }}

		{{ if gt ( len .PrimaryKeyFields ) 1 }}
			// sql query with composite primary key
			const sqlstr = `UPDATE {{ $table }} SET ` +
				`{{ colnamesquerymulti .Fields ", " 0 .PrimaryKeyFields }}` +
				` WHERE {{ colnamesquery .PrimaryKeyFields " AND " }}`

			// run query
			XOLog(sqlstr, {{ fieldnamesmulti .Fields $short .PrimaryKeyFields }}, {{ fieldnames .PrimaryKeyFields $short}})
			_, err = db.Exec(sqlstr, {{ fieldnamesmulti .Fields $short .PrimaryKeyFields }}, {{ fieldnames .PrimaryKeyFields $short}})
			return err
		{{- else }}
			// sql query
			const sqlstr = `UPDATE {{ $table }} SET ` +
				`{{ colnamesquery .Fields ", " .PrimaryKey.Name }}` +
				` WHERE {{ colname .PrimaryKey.Col }} = ?`

			// run query
			XOLog(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }}, {{ $short }}.{{ .PrimaryKey.Name }})
			_, err = db.Exec(sqlstr, {{ fieldnames .Fields $short .PrimaryKey.Name }}, {{ $short }}.{{ .PrimaryKey.Name }})
			return err
		{{- end }}
	}
{{ else }}
	// Update statements omitted due to lack of fields other than primary key
{{ end }}

// Delete deletes the {{ .Name }} from the database.
func ({{ $short }} *{{ .Name }}) Delete(db XODB) error {
	var err error

	{{ if gt ( len .PrimaryKeyFields ) 1 }}
		// sql query with composite primary key
		const sqlstr = `DELETE FROM {{ $table }} WHERE {{ colnamesquery .PrimaryKeyFields " AND " }}`

		// run query
		XOLog(sqlstr, {{ fieldnames .PrimaryKeyFields $short }})
		_, err = db.Exec(sqlstr, {{ fieldnames .PrimaryKeyFields $short }})
		if err != nil {
			return err
		}
	{{- else }}
		// sql query
		const sqlstr = `DELETE FROM {{ $table }} WHERE {{ colname .PrimaryKey.Col }} = ?`

		// run query
		XOLog(sqlstr, {{ $short }}.{{ .PrimaryKey.Name }})
		_, err = db.Exec(sqlstr, {{ $short }}.{{ .PrimaryKey.Name }})
		if err != nil {
			return err
		}
	{{- end }}

	return nil
}
{{- end }}

func {{ .Name }}Search(db XODB, where Where, order Order, limit, offset int) ([]*{{ .Name }}, error) {
	querybldr := sq.Select("*").From("{{ $table }}")

	for column, opAndVals := range where {
		if len(opAndVals) < 2 {
			return nil, errors.New("[" + column + "] invalid params")
		}
		var vals interface{}
		if len(opAndVals[1:]) == 1 {
			vals = opAndVals[1]
		} else {
			vals = opAndVals[1:]
		}
		switch opAndVals[0] {
		case "=":
			querybldr = querybldr.Where(sq.Eq{column: vals})
		case "!=":
			querybldr = querybldr.Where(sq.NotEq{column: vals})
		case "<":
			querybldr = querybldr.Where(sq.Gt{column: vals})
		case "<=":
			querybldr = querybldr.Where(sq.GtOrEq{column: vals})
		case ">":
			querybldr = querybldr.Where(sq.Lt{column: vals})
		case ">=":
			querybldr = querybldr.Where(sq.LtOrEq{column: vals})
		case "like":
			querybldr = querybldr.Where(sq.Like{column: vals})
		default:
			return nil, errors.New("invalid operator '"+ opAndVals[0] + "'")
		}
	}

	if len(order) > 0 {
		if len(order) != 2 {
			return nil, errors.New("invalid order")
		}
		ascOrDesc := strings.ToLower(order[1])
		if ascOrDesc != "asc" {
			if ascOrDesc != "desc" {
				return nil, errors.New("invalid order")
			}
		}
		querybldr = querybldr.OrderBy(strings.Join([]string(order), " "))
	}

	sqlstr, vals, _ := querybldr.Limit(uint64(limit)).Offset(uint64(offset)).ToSql()

	var res []*{{ .Name }}
	err := db.Select(&res, sqlstr, vals...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func {{ .Name }}SearchCount(db XODB, where Where) (int, error) {
	querybldr := sq.Select("count(*)").From("{{ $table }}")

	for column, opAndVals := range where {
		if len(opAndVals) < 2 {
			return 0, errors.New("[" + column + "] invalid params")
		}
		var vals interface{}
		if len(opAndVals[1:]) == 1 {
			vals = opAndVals[1]
		} else {
			vals = opAndVals[1:]
		}
		switch opAndVals[0] {
		case "=":
			querybldr = querybldr.Where(sq.Eq{column: vals})
		case "!=":
			querybldr = querybldr.Where(sq.NotEq{column: vals})
		case "<":
			querybldr = querybldr.Where(sq.Gt{column: vals})
		case "<=":
			querybldr = querybldr.Where(sq.GtOrEq{column: vals})
		case ">":
			querybldr = querybldr.Where(sq.Lt{column: vals})
		case ">=":
			querybldr = querybldr.Where(sq.LtOrEq{column: vals})
		case "like":
			querybldr = querybldr.Where(sq.Like{column: vals})
		default:
			return 0, errors.New("invalid operator '"+ opAndVals[0] + "'")
		}
	}

	sqlstr, vals, _ := querybldr.ToSql()

	var res int
	err := db.Get(&res, sqlstr, vals...)
	if err != nil {
		return 0, err
	}
	return res, nil
}
