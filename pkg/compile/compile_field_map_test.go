package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphemap-ts/pkg/compile"
	"github.com/kalo-build/plugin-morphemap-ts/pkg/mapdef"
)

type FieldMapTestSuite struct {
	suite.Suite
}

func TestFieldMapTestSuite(t *testing.T) {
	suite.Run(t, new(FieldMapTestSuite))
}

func (suite *FieldMapTestSuite) TestBuildTsFieldMapData_SimpleScalarMapping() {
	m := &mapdef.MorpheMap{
		Name: "OrgToProject",
		Aliases: map[string]string{
			"Org":  "Organization",
			"Proj": "Project",
		},
		Fields: mapdef.FieldMappings{
			"Proj.Name": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Org.Name"},
			"Proj.Code": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Org.Code"},
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
	}

	data := compile.BuildTsFieldMapData(m, config)

	suite.Equal("orgToProject", data.ConverterName)
	suite.Equal("Organization", data.SourceType)
	suite.Equal("Project", data.TargetType)
	suite.Equal("@/types/source", data.SourceTypesImportPath)
	suite.Equal("@/types/target", data.TargetTypesImportPath)
	suite.Len(data.Fields, 2)
}

func (suite *FieldMapTestSuite) TestBuildTsFieldMapData_ConstantValue() {
	m := &mapdef.MorpheMap{
		Name: "WithConstant",
		Aliases: map[string]string{
			"Src": "Source",
			"Tgt": "Target",
		},
		Fields: mapdef.FieldMappings{
			"Tgt.Status": mapdef.FieldMappingValue{IsScalar: true, Scalar: "active"},
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
	}

	data := compile.BuildTsFieldMapData(m, config)

	suite.Require().Len(data.Fields, 1)
	suite.True(data.Fields[0].IsConstant)
	suite.Equal("'active'", data.Fields[0].ConstValue)
}

func (suite *FieldMapTestSuite) TestBuildTsFieldMapData_ObjectWithCast() {
	m := &mapdef.MorpheMap{
		Name: "WithCast",
		Aliases: map[string]string{
			"Src": "Source",
			"Tgt": "Target",
		},
		Fields: mapdef.FieldMappings{
			"Tgt.Count": mapdef.FieldMappingValue{
				IsScalar: false,
				Object: &mapdef.FieldMapping{
					From: "Src.Count",
					Cast: "Number",
				},
			},
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
	}

	data := compile.BuildTsFieldMapData(m, config)

	suite.Require().Len(data.Fields, 1)
	suite.Equal("Number", data.Fields[0].Cast)
}

func (suite *FieldMapTestSuite) TestBuildTsFieldMapData_ObjectWithValueMap() {
	m := &mapdef.MorpheMap{
		Name: "WithValueMap",
		Aliases: map[string]string{
			"Src": "Source",
			"Tgt": "Target",
		},
		Fields: mapdef.FieldMappings{
			"Tgt.Status": mapdef.FieldMappingValue{
				IsScalar: false,
				Object: &mapdef.FieldMapping{
					From: "Src.Status",
					ValueMap: map[string]string{
						"active":   "Active",
						"inactive": "Inactive",
					},
				},
			},
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
	}

	data := compile.BuildTsFieldMapData(m, config)

	suite.Require().Len(data.Fields, 1)
	suite.True(data.Fields[0].HasValueMap)
	suite.Len(data.Fields[0].ValueMap, 2)
}

func (suite *FieldMapTestSuite) TestBuildTsFieldMapData_WithHooks() {
	m := &mapdef.MorpheMap{
		Name: "WithHooks",
		Aliases: map[string]string{
			"Src": "Source",
			"Tgt": "Target",
		},
		Fields: mapdef.FieldMappings{
			"Tgt.Name": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Src.Name"},
		},
		Hooks: &mapdef.Hooks{
			AfterMap: []mapdef.Hook{
				{Name: "applyDefaults"},
			},
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
	}

	data := compile.BuildTsFieldMapData(m, config)

	suite.Require().Len(data.Hooks, 1)
	suite.Equal("applyDefaults", data.Hooks[0])
}

func (suite *FieldMapTestSuite) TestBuildTsFieldMapData_CamelCaseConversion() {
	m := &mapdef.MorpheMap{
		Name: "CamelCaseTest",
		Aliases: map[string]string{
			"Src": "Source",
			"Tgt": "Target",
		},
		Fields: mapdef.FieldMappings{
			"Tgt.FirstName": mapdef.FieldMappingValue{IsScalar: true, Scalar: "Src.FirstName"},
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
	}

	data := compile.BuildTsFieldMapData(m, config)

	suite.Require().Len(data.Fields, 1)
	// PascalCase fields should be converted to camelCase in TS
	suite.Equal("firstName", data.Fields[0].TargetField)
	suite.Equal("firstName", data.Fields[0].SourceExpr)
}

func (suite *FieldMapTestSuite) TestRenderTsFieldMapTemplate_ProducesValidTs() {
	data := compile.TsFieldMapData{
		ConverterName:         "orgToProject",
		SourceType:            "Organization",
		TargetType:            "Project",
		SourceTypesImportPath: "@/types/models",
		TargetTypesImportPath: "@/types/models",
		Fields: []compile.TsFieldMappingData{
			{TargetField: "name", SourceExpr: "name"},
			{TargetField: "code", SourceExpr: "code"},
		},
	}

	result, renderErr := compile.RenderTsFieldMapTemplate(data)

	suite.NoError(renderErr)
	suite.Contains(result, "import type { Organization } from '@/types/models'")
	suite.Contains(result, "import type { Project } from '@/types/models'")
	suite.Contains(result, "export function orgToProject(source: Organization): Project")
	suite.Contains(result, "target.name = source.name")
	suite.Contains(result, "target.code = source.code")
	suite.Contains(result, "return target as Project")
}

func (suite *FieldMapTestSuite) TestRenderTsFieldMapTemplate_ConstantAssignment() {
	data := compile.TsFieldMapData{
		ConverterName:         "testConv",
		SourceType:            "Src",
		TargetType:            "Tgt",
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
		Fields: []compile.TsFieldMappingData{
			{TargetField: "status", IsConstant: true, ConstValue: "'active'"},
		},
	}

	result, renderErr := compile.RenderTsFieldMapTemplate(data)

	suite.NoError(renderErr)
	suite.Contains(result, "target.status = 'active'")
}

func (suite *FieldMapTestSuite) TestRenderTsFieldMapTemplate_CastAssignment() {
	data := compile.TsFieldMapData{
		ConverterName:         "testConv",
		SourceType:            "Src",
		TargetType:            "Tgt",
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
		Fields: []compile.TsFieldMappingData{
			{TargetField: "count", SourceExpr: "count", Cast: "Number"},
		},
	}

	result, renderErr := compile.RenderTsFieldMapTemplate(data)

	suite.NoError(renderErr)
	suite.Contains(result, "target.count = Number(source.count) as any")
}
