package BusinessLogicModels

type EmployeeDTO struct {
	Name                  string   `json:"name"`
	UsableMachineNameList []string `json:"usable_machine_name_list"`
}
