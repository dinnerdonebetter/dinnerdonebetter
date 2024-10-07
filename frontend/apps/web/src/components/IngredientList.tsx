import { List, Checkbox } from '@mantine/core';

import { Recipe, RecipeStepIngredient } from '@dinnerdonebetter/models';
import { cleanFloat, determineAllIngredientsForRecipes } from '@dinnerdonebetter/utils';

declare interface IngredientListComponentProps {
  scale: number;
  recipes: Recipe[];
}

export const RecipeIngredientListComponent = ({ recipes, scale }: IngredientListComponentProps): JSX.Element => {
  recipes.forEach((recipe: Recipe) => {
    recipes = recipes.concat(recipe.supportingRecipes || []);
  });

  return (
    <List icon={<></>} pb="sm">
      {determineAllIngredientsForRecipes(
        recipes.map((x) => {
          return { recipe: x, scale };
        }),
      ).map((ingredient: RecipeStepIngredient) => {
        let measurmentUnitName =
          ingredient.quantity.min === 1 ? ingredient.measurementUnit.name : ingredient.measurementUnit.pluralName;

        let minQty = cleanFloat(scale === 1.0 ? ingredient.quantity.min : ingredient.quantity.min * scale);

        let maxQty = cleanFloat(
          scale === 1.0 ? (ingredient.quantity.max ?? 0) : (ingredient.quantity.max ?? 0 * scale),
        );

        return (
          <List.Item key={ingredient.id} my="-sm">
            <Checkbox
              label={
                <>
                  <u>
                    {` ${minQty}${maxQty > 0 ? `- ${maxQty}` : ''} ${
                      ['unit', 'units'].includes(measurmentUnitName) ? '' : measurmentUnitName
                    }`}
                  </u>{' '}
                  {ingredient.name}
                </>
              }
            />
          </List.Item>
        );
      })}
    </List>
  );
};
