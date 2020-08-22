package bleve

import (
	bleve "github.com/blevesearch/bleve"
	en "github.com/blevesearch/bleve/analysis/lang/en"
	mapping "github.com/blevesearch/bleve/mapping"
)

func buildValidInstrumentMapping() *mapping.IndexMappingImpl {
	m := mapping.NewIndexMapping()

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	validInstrumentMapping := bleve.NewDocumentMapping()
	validInstrumentMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("variant", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("icon", englishTextFieldMapping)
	validInstrumentMapping.AddFieldMappingsAt("belongsToUser", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("validInstrument", validInstrumentMapping)

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
	validIngredientMapping.AddFieldMappingsAt("icon", englishTextFieldMapping)
	validIngredientMapping.AddFieldMappingsAt("belongsToUser", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("validIngredient", validIngredientMapping)

	return m
}

func buildValidPreparationMapping() *mapping.IndexMappingImpl {
	m := mapping.NewIndexMapping()

	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	validPreparationMapping := bleve.NewDocumentMapping()
	validPreparationMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	validPreparationMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	validPreparationMapping.AddFieldMappingsAt("icon", englishTextFieldMapping)
	validPreparationMapping.AddFieldMappingsAt("belongsToUser", bleve.NewNumericFieldMapping())
	m.AddDocumentMapping("validPreparation", validPreparationMapping)

	return m
}
