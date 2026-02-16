package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphemap-ts/pkg/compile"
	"github.com/kalo-build/plugin-morphemap-ts/pkg/mapdef"
)

type EnumMapTestSuite struct {
	suite.Suite
}

func TestEnumMapTestSuite(t *testing.T) {
	suite.Run(t, new(EnumMapTestSuite))
}

func (suite *EnumMapTestSuite) TestBuildTsEnumMapData_BasicEntries() {
	m := &mapdef.MorpheMap{
		Name: "StatusToPriority",
		Aliases: map[string]string{
			"Src": "Status",
			"Tgt": "Priority",
		},
		Entries: map[string]string{
			"Tgt.Low":    "Src.Inactive",
			"Tgt.Medium": "Src.Pending",
			"Tgt.High":   "Src.Active",
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/enums",
		TargetTypesImportPath: "@/types/enums",
	}

	data := compile.BuildTsEnumMapData(m, config)

	suite.Equal("statusToPriority", data.FunctionName)
	suite.Equal("Status", data.SourceType)
	suite.Equal("Priority", data.TargetType)
	suite.Len(data.Entries, 3)
}

func (suite *EnumMapTestSuite) TestBuildTsEnumMapData_StripsAliasPrefix() {
	m := &mapdef.MorpheMap{
		Name: "SimpleEnum",
		Aliases: map[string]string{
			"A": "EnumA",
			"B": "EnumB",
		},
		Entries: map[string]string{
			"B.ValueX": "A.ValueY",
		},
	}
	config := compile.TsConverterConfig{
		SourceTypesImportPath: "@/types/source",
		TargetTypesImportPath: "@/types/target",
	}

	data := compile.BuildTsEnumMapData(m, config)

	suite.Require().Len(data.Entries, 1)
	suite.Equal("ValueX", data.Entries[0].TargetEntry)
	suite.Equal("ValueY", data.Entries[0].SourceEntry)
}

func (suite *EnumMapTestSuite) TestBuildTsEnumMapData_CamelCaseFunctionName() {
	m := &mapdef.MorpheMap{
		Name: "MySourceToMyTarget",
		Aliases: map[string]string{
			"Src": "EnumA",
			"Tgt": "EnumB",
		},
		Entries: map[string]string{
			"Tgt.A": "Src.B",
		},
	}
	config := compile.TsConverterConfig{}

	data := compile.BuildTsEnumMapData(m, config)

	suite.Equal("mySourceToMyTarget", data.FunctionName)
}

func (suite *EnumMapTestSuite) TestRenderTsEnumMapTemplate_ProducesValidTs() {
	data := compile.TsEnumMapData{
		FunctionName: "statusToPriority",
		SourceType:   "Status",
		TargetType:   "Priority",
		Entries: []compile.TsEnumEntryData{
			{SourceEntry: "Active", TargetEntry: "High"},
			{SourceEntry: "Inactive", TargetEntry: "Low"},
		},
	}

	result, renderErr := compile.RenderTsEnumMapTemplate(data)

	suite.NoError(renderErr)
	suite.Contains(result, "const statusToPriorityMap: Record<string, string>")
	suite.Contains(result, "'Active': 'High'")
	suite.Contains(result, "'Inactive': 'Low'")
	suite.Contains(result, "export function statusToPriority(source: string): string")
	suite.Contains(result, "throw new Error")
}

func (suite *EnumMapTestSuite) TestRenderTsEnumMapTemplate_EmptyEntries() {
	data := compile.TsEnumMapData{
		FunctionName: "emptyEnum",
		SourceType:   "TypeA",
		TargetType:   "TypeB",
		Entries:      []compile.TsEnumEntryData{},
	}

	result, renderErr := compile.RenderTsEnumMapTemplate(data)

	suite.NoError(renderErr)
	suite.Contains(result, "export function emptyEnum(source: string): string")
}
