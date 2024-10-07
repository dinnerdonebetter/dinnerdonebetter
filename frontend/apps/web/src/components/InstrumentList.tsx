import { List, Checkbox } from '@mantine/core';

import { Recipe, RecipeStepInstrument, RecipeStepVessel } from '@dinnerdonebetter/models';
import { determineAllInstrumentsForRecipes } from '@dinnerdonebetter/utils';

declare interface InstrumentListComponentProps {
  recipes: Recipe[];
}

export const RecipeInstrumentListComponent = ({ recipes }: InstrumentListComponentProps): JSX.Element => {
  recipes.forEach((recipe: Recipe) => {
    recipes = recipes.concat(recipe.supportingRecipes || []);
  });

  return (
    <List icon={<></>} pb="sm">
      {determineAllInstrumentsForRecipes(recipes).map((x: RecipeStepInstrument | RecipeStepVessel) => {
        const isRecipeStepInstrument = JSON.parse(JSON.stringify(x)).hasOwnProperty('instrument');

        if (isRecipeStepInstrument) {
          return (
            <List.Item key={x.id} my="-sm">
              <Checkbox
                size="sm"
                label={`${x.quantity.min}${
                  (x.quantity.max ?? 0) > 0 && x.quantity.max != x.quantity.min ? `- ${x.quantity.max}` : ''
                } ${(x as RecipeStepInstrument).instrument?.name}`}
              />
            </List.Item>
          );
        } else {
          return (
            <List.Item key={x.id} my="-sm">
              <Checkbox
                size="sm"
                label={`${x.quantity.min}${
                  (x.quantity.max ?? 0) > 0 && x.quantity.max != x.quantity.min ? `- ${x.quantity.max}` : ''
                } ${(x as RecipeStepVessel).vessel?.name}`}
              />
            </List.Item>
          );
        }
      })}
    </List>
  );
};
