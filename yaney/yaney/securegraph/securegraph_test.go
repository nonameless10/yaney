package securegraph

import (
	"github.com/nonameless10/yaney/securekv/sse"
	"testing"
)

func getSecret(passphrase, salt string, iter int) []byte {
	return sse.Key([]byte(passphrase), []byte(salt), iter)
}

//func getVertexGraph() []VertexGraph {
//	vgs := make([]VertexGraph, 0)
//
//	// 乔峰
//	var vg1 VertexGraph
//	vg1.vertex = Vertex{
//		name: "乔峰",
//		props: map[string]string{
//			"类型": "人物",
//			"年龄": "36",
//			"职务": "丐帮帮主",
//		},
//	}
//	vg1.outEdges = append(vg1.outEdges, Edge{
//		selfVertexName: "乔峰",
//		inout:          1,
//		relationName:   "会",
//		rank:           0,
//		peerVertexName: "降龙十八掌",
//		props: map[string]string{
//			"掌握时长": "10年",
//			"掌握程度": "精通",
//		},
//	})
//	vg1.outEdges = append(vg1.outEdges, Edge{
//		selfVertexName: "乔峰",
//		inout:          1,
//		relationName:   "结拜",
//		rank:           0,
//		peerVertexName: "段誉",
//		props: map[string]string{
//			"结拜时间": "金历233年冬",
//			"结拜时长": "5年",
//		},
//	})
//	vg1.outEdges = append(vg1.outEdges, Edge{
//		selfVertexName: "乔峰",
//		inout:          1,
//		relationName:   "结拜",
//		rank:           0,
//		peerVertexName: "虚竹",
//		props: map[string]string{
//			"结拜时间": "金历233年冬",
//			"结拜时长": "5年",
//		},
//	})
//
//	// 段誉
//	var vg2 VertexGraph
//	vg2.vertex = Vertex{
//		name: "段誉",
//		props: map[string]string{
//			"类型": "人物",
//			"年龄": "27",
//			"职务": "大理王爷",
//		},
//	}
//	vg2.outEdges = append(vg2.outEdges, Edge{
//		selfVertexName: "段誉",
//		inout:          1,
//		relationName:   "会",
//		rank:           0,
//		peerVertexName: "凌波微步",
//		props: map[string]string{
//			"掌握时长": "3年",
//			"掌握程度": "熟悉",
//		},
//	})
//
//	vg2.outEdges = append(vg2.outEdges, Edge{
//		selfVertexName: "段誉",
//		inout:          1,
//		relationName:   "结拜",
//		rank:           0,
//		peerVertexName: "乔峰",
//		props: map[string]string{
//			"结拜时间": "金历233年冬",
//			"结拜时长": "5年",
//		},
//	})
//	vg2.outEdges = append(vg2.outEdges, Edge{
//		selfVertexName: "段誉",
//		inout:          1,
//		relationName:   "结拜",
//		rank:           0,
//		peerVertexName: "虚竹",
//		props: map[string]string{
//			"结拜时间": "金历233年冬",
//			"结拜时长": "5年",
//		},
//	})
//
//	// 虚竹
//	var vg3 VertexGraph
//	vg3.vertex = Vertex{
//		name: "虚竹",
//		props: map[string]string{
//			"类型": "人物",
//			"年龄": "24",
//			"职务": "逍遥派掌门",
//		},
//	}
//	vg3.outEdges = append(vg3.outEdges, Edge{
//		selfVertexName: "虚竹",
//		inout:          1,
//		relationName:   "会",
//		rank:           0,
//		peerVertexName: "小无相功",
//		props: map[string]string{
//			"掌握时长": "1年",
//			"掌握程度": "略懂",
//		},
//	})
//	vg3.outEdges = append(vg3.outEdges, Edge{
//		selfVertexName: "虚竹",
//		inout:          1,
//		relationName:   "结拜",
//		rank:           0,
//		peerVertexName: "乔峰",
//		props: map[string]string{
//			"结拜时间": "金历233年冬",
//			"结拜时长": "5年",
//		},
//	})
//	vg3.outEdges = append(vg3.outEdges, Edge{
//		selfVertexName: "虚竹",
//		inout:          1,
//		relationName:   "结拜",
//		rank:           0,
//		peerVertexName: "段誉",
//		props: map[string]string{
//			"结拜时间": "金历233年冬",
//			"结拜时长": "5年",
//		},
//	})
//
//	// 武功
//	var vg4 VertexGraph
//	vg4.vertex = Vertex{
//		name: "降龙十八掌",
//		props: map[string]string{
//			"类型": "外功",
//			"强度": "5",
//		},
//	}
//	vg4.inEdges = append(vg4.inEdges, Edge{
//		selfVertexName: "降龙十八掌",
//		inout:          -1,
//		relationName:   "会",
//		rank:           0,
//		peerVertexName: "乔峰",
//		props: map[string]string{
//			"掌握时长": "10年",
//			"掌握程度": "精通",
//		},
//	})
//
//	var vg5 VertexGraph
//	vg5.vertex = Vertex{
//		name: "凌波微步",
//		props: map[string]string{
//			"类型": "轻功",
//			"强度": "6",
//		},
//	}
//	vg5.inEdges = append(vg5.inEdges, Edge{
//		selfVertexName: "凌波微步",
//		inout:          -1,
//		relationName:   "会",
//		rank:           0,
//		peerVertexName: "段誉",
//		props: map[string]string{
//			"掌握时长": "3年",
//			"掌握程度": "熟悉",
//		},
//	})
//
//	var vg6 VertexGraph
//	vg6.vertex = Vertex{
//		name: "小无相功",
//		props: map[string]string{
//			"类型": "内功",
//			"强度": "6",
//		},
//	}
//	vg6.inEdges = append(vg6.inEdges, Edge{
//		selfVertexName: "小无相功",
//		inout:          -1,
//		relationName:   "会",
//		rank:           0,
//		peerVertexName: "虚竹",
//		props: map[string]string{
//			"掌握时长": "1年",
//			"掌握程度": "略懂",
//		},
//	})
//
//	vgs = append(vgs, vg1)
//	vgs = append(vgs, vg2)
//	vgs = append(vgs, vg3)
//	vgs = append(vgs, vg4)
//	vgs = append(vgs, vg5)
//	vgs = append(vgs, vg6)
//
//	return vgs
//}

func TestNewSecureGraph(t *testing.T) {
	masterSecret := getSecret("master", "secret", 4096)
	keySecret := getSecret("key", "secret", 4096)
	valSecret := getSecret("val", "secret", 4096)

	sg := NewSecureGraph("graphdir", masterSecret, keySecret, valSecret)
	sg.BuildGraph(getVertexGraph())

	vg1, _ := sg.QueryGraph("乔峰")
	t.Log(vg1)

	vg2, _ := sg.QueryGraph("降龙十八掌")
	t.Log(vg2)

}
