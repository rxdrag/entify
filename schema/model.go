package schema

import (
	"github.com/mitchellh/mapstructure"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"

	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/utils"
)

func QueryPublishedMeta() interface{} {
	publishedMeta := repository.QueryOneEntity(model.GlobalModel.Graph.GetMetaEntity(), graph.QueryArg{
		consts.ARG_WHERE: graph.QueryArg{
			consts.META_STATUS: graph.QueryArg{
				consts.ARG_EQ: meta.META_STATUS_PUBLISHED,
			},
		},
	}, repository.NewSupperVerifier())

	return publishedMeta
}

func QueryNextMeta() interface{} {
	nextMeta := repository.QueryOneEntity(model.GlobalModel.Graph.GetMetaEntity(), graph.QueryArg{
		consts.ARG_WHERE: graph.QueryArg{
			consts.META_STATUS: graph.QueryArg{
				consts.ARG_ISNULL: true,
			},
		},
	}, repository.NewSupperVerifier())

	return nextMeta
}

func DecodeContent(obj interface{}) *meta.MetaContent {
	content := meta.MetaContent{}
	if obj != nil {
		err := mapstructure.Decode(obj.(utils.Object)[consts.META_CONTENT], &content)
		if err != nil {
			panic("Decode content failure:" + err.Error())
		}
	}
	return &content
}

func LoadModel() {
	//初始值，用户取meta信息，取完后，换掉该部分内容
	initMeta := meta.MetaContent{
		Classes: []meta.ClassMeta{
			meta.MetaStatusEnum,
			meta.MetaClass,
			meta.EntityAuthSettingsClass,
			meta.AbilityTypeEnum,
			meta.AbilityClass,
		},
	}
	model.GlobalModel = model.New(&initMeta)
	publishedMeta := QueryPublishedMeta()
	publishedContent := DecodeContent(publishedMeta)
	publishedContent.Classes = append(publishedContent.Classes,
		meta.MetaStatusEnum,
		meta.MetaClass,
		meta.EntityAuthSettingsClass,
		meta.AbilityTypeEnum,
		meta.AbilityClass,
	)

	model.GlobalModel = model.New(publishedContent)
}

// func init() {
// 	LoadModel()
// }
