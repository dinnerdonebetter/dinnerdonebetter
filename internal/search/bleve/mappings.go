package bleve

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
)

func buildValidInstrumentMapping() *mapping.IndexMappingImpl {
	m := mapping.NewIndexMapping()

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	validInstrumentMapping := bleve.NewDocumentMapping()
	validInstrumentMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("variant", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("iconPath", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("belongsToHousehold", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("valid_instrument", validInstrumentMapping)

	return m
}

func buildValidPreparationMapping() *mapping.IndexMappingImpl {
	m := mapping.NewIndexMapping()

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	validPreparationMapping := bleve.NewDocumentMapping()
	validPreparationMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	validPreparationMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	validPreparationMapping.AddFieldMappingsAt("iconPath", englishTextFieldMapping)
	validPreparationMapping.AddFieldMappingsAt("belongsToHousehold", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("valid_preparation", validPreparationMapping)

	return m
}

func buildValidIngredientMapping() *mapping.IndexMappingImpl {
	m := mapping.NewIndexMapping()

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	validIngredientMapping := bleve.NewDocumentMapping()
	validIngredientMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	validIngredientMapping.AddFieldMappingsAt("variant", englishTextFieldMapping)
	validIngredientMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	validIngredientMapping.AddFieldMappingsAt("warning", englishTextFieldMapping)
	validIngredientMapping.AddFieldMappingsAt("iconPath", englishTextFieldMapping)
	validIngredientMapping.AddFieldMappingsAt("belongsToHousehold", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("valid_ingredient", validIngredientMapping)

	return m
}
