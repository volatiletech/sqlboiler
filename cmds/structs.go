package cmds

// init the "struct" command
// func init() {
// 	SQLBoiler.AddCommand(structCmd)
// 	structCmd.Run = structRun
// }

// var structCmd = &cobra.Command{
// 	Use:   "struct",
// 	Short: "Generate structs from table definitions",
// }
//
// // deleteRun executes the struct command, and generates the struct definitions
// // boilerplate from the template file.
// func structRun(cmd *cobra.Command, args []string) {
// 	err := outHandler(generateStructs())
// 	if err != nil {
// 		errorQuit(err)
// 	}
// }

// generateStructs returns a slice of each template execution result.
// Each of these results holds a struct definition generated from the struct template.
// func generateStructs() [][]byte {
// 	t, err := template.New("struct.tpl").Funcs(template.FuncMap{
// 		"makeGoColName": makeGoColName,
// 		"makeDBColName": makeDBColName,
// 	}).ParseFiles("templates/struct.tpl")
//
// 	if err != nil {
// 		errorQuit(err)
// 	}
//
// 	outputs, err := processTemplate(t)
// 	if err != nil {
// 		errorQuit(err)
// 	}
//
// 	return outputs
// }
