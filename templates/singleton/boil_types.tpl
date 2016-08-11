// M type is for providing columns and column values to UpdateAll.
type M map[string]interface{}

type upsertData struct {
  conflict  []string
  update    []string
  whitelist []string
  returning []string
}
