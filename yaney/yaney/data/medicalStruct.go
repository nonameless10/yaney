package data

type Disease struct {
	DiseaseNo     string `json:"diseaseNo"`
	DiagnosisDate string `json:"diagnosisDate"`
}

type Assay struct {
	AssayID   string `json:"assayID"`
	AssayDate string `json:"assayDate"`
	AssayRes  string `json:"assayRes"`
	AssayFlag string `json:"assayFlag"`
}
type Operation struct {
	OperationID         string `json:"operationID"`
	OperationDoctorName string `json:"operationDoctorName"`
	OperationDate       string `json:"operationDate"`
}
type BeHospitalMsg struct {
	HosID                    string `json:"hosID"`
	AdmissionDepartmentName  string `json:"admissionDepartmentName"`
	AdmissionBedNo           string `json:"admissionBedNo"`
	AdmissionDate            string `json:"admissionDate"`
	DischargedDepartmentName string `json:"dischargedDepartmentName"`
	DischargedBedNo          string `json:"dischargedBedNo"`
	DischargedDate           string `json:"dischargedDate"`
	DoctorAdvice             string `json:"doctorAdvice"`
}

type PatientData struct {
	PatientID         string          `json:"patientID"`
	Name              string          `json:"name"`
	Sex               string          `json:"sex"`
	Birthday          string          `json:"birthday"`
	Profession        string          `json:"profession"`
	AllergyMedication string          `json:"allergyMedication"`
	Diseases          []Disease       `json:"diseases"`
	Assays            []Assay         `json:"assays"`
	Operations        []Operation     `json:"operations"`
	BeHospitalMsgs    []BeHospitalMsg `json:"beHospitalMsgs"`
}

type AssaysData struct {
	AssayID   string `json:"assayID"`
	AssayName string `json:"assayName"`
	AssayUnit string `json:"assayUnit"`
}

type DiseaseData struct {
	DiseaseNo   string `json:"diseaseNo"`
	DiseaseName string `json:"diseaseName"`
}

type OperationData struct {
	OperationID   string `json:"operationID"`
	OperationName string `json:"operationName"`
}

type MedicalData struct {
	PatientDatas   []PatientData   `json:"patientData"`
	AssaysDatas    []AssaysData    `json:"assaysDatas"`
	DiseaseDatas   []DiseaseData   `json:"diseaseDatas"`
	OperationDatas []OperationData `json:"operationDatas"`
}

