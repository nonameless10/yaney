package securegraph

import (
	"encoding/json"
	"fmt"
	"github.com/jo3yzhu/yaney/data"
	"github.com/jo3yzhu/yaney/securekv/sse"
	"io/ioutil"
)

type VertexPair struct {
	SelfName string
	PeerName string
}



func GetSecret(passphrase, salt string, iter int) []byte {
	return sse.Key([]byte(passphrase), []byte(salt), iter)
}

// 添加住院信息实体和病人实体
func addPatientAndHospitalVertex(records map[string]*VertexGraph, datas *data.MedicalData){
	for _, patient :=range datas.PatientDatas{
		records[patient.PatientID] = new(VertexGraph)
		records[patient.PatientID].vertex = Vertex{
			name: patient.PatientID,
			props: map[string]string{
				"姓名":patient.Name,
				"性别":patient.Sex,
				"出生日期":patient.Birthday,
				"职位名称":patient.Profession,
				"过敏药物":patient.AllergyMedication,
			},
		}

		// 住院信息
		for _, hos := range patient.BeHospitalMsgs{
			records[hos.HosID] = new(VertexGraph)
			records[hos.HosID].vertex = Vertex{
				name: hos.HosID,
				props: map[string]string{
					"入院科室": hos.DischargedDepartmentName,
					"入院床号": hos.AdmissionBedNo,
					"入院时间": hos.AdmissionDate,
					"出院科室": hos.DischargedDepartmentName,
					"出院床号":  hos.DischargedBedNo,
					"出院时间": hos.DischargedDate,
					"医嘱": hos.DoctorAdvice,
				},
			}
		}
	}
}


// 添加化验项目实体
func addAssayVertex(records map[string]*VertexGraph, datas *data.MedicalData){
	for _, assay := range datas.AssaysDatas{
		records[assay.AssayID] = new(VertexGraph)
		records[assay.AssayID].vertex = Vertex{
			name: assay.AssayID,
			props: map[string]string{
				"化验项目名称":assay.AssayName,
				"单位": assay.AssayUnit,
			},
		}
	}
}

func addOperationVertex(records map[string]*VertexGraph, datas *data.MedicalData){
	for _, opt := range datas.OperationDatas{
		records[opt.OperationID] = new(VertexGraph)
		records[opt.OperationID].vertex = Vertex{
			name: opt.OperationID,
			props: map[string]string{
				"手术名称": opt.OperationName,
			},
		}
	}
}

func addDiseaseVertex(records map[string]*VertexGraph, datas *data.MedicalData){
	for _, disease := range  datas.DiseaseDatas{
		records[disease.DiseaseNo] = new(VertexGraph)
		records[disease.DiseaseNo].vertex = Vertex{
			name: disease.DiseaseNo,
			props: map[string]string{
				"疾病名称": disease.DiseaseName,
			},
		}
	}
}


func dealEdge(records map[string]*VertexGraph, datas *data.MedicalData, rankRecord map[VertexPair]uint32){
	for _, patient := range datas.PatientDatas{
		// 1.遍历住院边
		for _, hos := range patient.BeHospitalMsgs{
			vp := VertexPair{patient.PatientID, hos.HosID}
			if _,ok := rankRecord[vp];!ok{
				rankRecord[vp] = 0
			}else{
				rankRecord[vp] = rankRecord[vp]+1
			}
			records[patient.PatientID].outEdges = append(records[patient.PatientID].outEdges, Edge{
				selfVertexName: patient.PatientID,
				inout: 1,
				relationName: "住院",
				rank: rankRecord[vp],
				peerVertexName:  hos.HosID,
				props: map[string]string{},
			})

			records[hos.HosID].inEdges = append(records[hos.HosID].inEdges, Edge{
				selfVertexName: hos.HosID,
				inout: -1,
				relationName: "住院",
				rank: rankRecord[vp],
				peerVertexName:  patient.PatientID,
				props: map[string]string{},
			})
		}
		// 2.遍历疾病边
		for _,disease := range patient.Diseases{
			vp := VertexPair{patient.PatientID, disease.DiseaseNo}
			if _,ok := rankRecord[vp];!ok{
				rankRecord[vp] = 0
			}else{
				rankRecord[vp] = rankRecord[vp]+1
			}
			records[patient.PatientID].outEdges = append(records[patient.PatientID].outEdges, Edge{
				selfVertexName: patient.PatientID,
				inout: 1,
				relationName: "患有",
				rank: rankRecord[vp],
				peerVertexName:  disease.DiseaseNo,
				props: map[string]string{
					"确诊时间": disease.DiagnosisDate,
				},
			})

			records[disease.DiseaseNo].inEdges = append(records[disease.DiseaseNo].inEdges, Edge{
				selfVertexName: disease.DiseaseNo,
				inout: -1,
				relationName: "患有",
				rank: rankRecord[vp],
				peerVertexName:  patient.PatientID,
				props: map[string]string{
					"确诊时间": disease.DiagnosisDate,
				},
			})
		}

		// 3.化验项目
		for _, assay := range patient.Assays{
			vp := VertexPair{patient.PatientID, assay.AssayID}
			if _,ok := rankRecord[vp];!ok{
				rankRecord[vp] = 0
			}else{
				rankRecord[vp] = rankRecord[vp]+1
			}
			records[patient.PatientID].outEdges = append(records[patient.PatientID].outEdges, Edge{
				selfVertexName: patient.PatientID,
				inout: 1,
				relationName: "化验",
				rank: rankRecord[vp],
				peerVertexName:  assay.AssayID,
				props: map[string]string{
					"化验日期": assay.AssayDate,
					"化验结果": assay.AssayRes,
					"化验结果标志": assay.AssayFlag,
				},
			})

			records[assay.AssayID].inEdges = append(records[assay.AssayID].inEdges, Edge{
				selfVertexName: assay.AssayID,
				inout: -1,
				relationName: "化验",
				rank: rankRecord[vp],
				peerVertexName:  patient.PatientID,
				props: map[string]string{
					"化验日期": assay.AssayDate,
					"化验结果": assay.AssayRes,
					"化验结果标志": assay.AssayFlag,
				},
			})
		}
		// 手术
		for _, opt := range patient.Operations{
			vp := VertexPair{patient.PatientID, opt.OperationID}
			if _,ok := rankRecord[vp];!ok{
				rankRecord[vp] = 0
			}else{
				rankRecord[vp] = rankRecord[vp]+1
			}

			records[patient.PatientID].outEdges = append(records[patient.PatientID].outEdges, Edge{
				selfVertexName: patient.PatientID,
				inout: 1,
				relationName: "手术",
				rank: rankRecord[vp],
				peerVertexName:  opt.OperationID,
				props: map[string]string{
					"手术医生":   opt.OperationDoctorName,
					"手术日期": opt.OperationDate,
				},
			})

			records[opt.OperationID].inEdges = append(records[opt.OperationID].inEdges, Edge{
				selfVertexName: opt.OperationID,
				inout: -1,
				relationName: "手术",
				rank: rankRecord[vp],
				peerVertexName:  patient.PatientID,
				props: map[string]string{
					"手术医生":   opt.OperationDoctorName,
					"手术日期": opt.OperationDate,
				},
			})
		}

	}
}



func getVertexGraph() []VertexGraph {
	var (
		err      error
		content  []byte
		datas     data.MedicalData
		filename string
		vgs []VertexGraph
	)
	vgs = make([]VertexGraph, 0)

	filename = "/home/l1nkkk/project/mime/yaney/data/data.json"

	if content, err = ioutil.ReadFile(filename); err != nil {
		return vgs
	}
	if err = json.Unmarshal(content, &datas); err != nil {
		return vgs
	}
	//fmt.Printf("%v",data)

	// string为实体id，每个实体有一个一度子图
	records := make(map[string]*VertexGraph)
	addPatientAndHospitalVertex(records, &datas)
	addAssayVertex(records, &datas)
	addDiseaseVertex(records, &datas)
	addOperationVertex(records, &datas)

	rankRecord := make(map[VertexPair]uint32)
	dealEdge(records, &datas, rankRecord)

	for _,vg := range records{
		vgs = append(vgs, *vg)
	}

	return vgs
}
func CreateMedicalGraph(){
	masterSecret := GetSecret("master", "secret", 4096)
	keySecret := GetSecret("key", "secret", 4096)
	valSecret := GetSecret("val", "secret", 4096)
	sg := NewSecureGraph("graphdir", masterSecret, keySecret, valSecret)
	sg.BuildGraph(getVertexGraph())
	vg1, _ := sg.QueryGraph("1908035220")
	fmt.Println(vg1)
}