{{- if .Table.IsJoinTable -}}
{{- else if false -}}
func TestIt(t *testing.T) {
  var a A
  var b, c B

  if err := a.Insert(); err != nil {
    t.Fatal(err)
  }

  b.user_id, c.user_id = a.ID, a.ID
  if err := b.Insert(); err != nil {
    t.Fatal(err)
  }
  if err := c.Insert(); err != nil {
    t.Fatal(err)
  }
}

{{- end -}}{{- /* outer if join table */ -}}
