<template>
  <div class="grid gap-4 grid-cols-1 justify-items-center mt-3">
    <div class="w-1/2 p-6 rounded-lg shadow-lg" :class="mealPlan.completed ? 'bg-gray-200' : 'bg-gray-50'" v-for="mealPlan in mealPlans">
      <h2 class="text-2xl font-bold mb-2 text-gray-800">{{ mealPlan.startDate }} Meal Plan</h2>
    </div>
  </div>
</template>


<script lang="ts">
import { defineComponent} from "vue";

class MealPlan {
  startDate: string;
  completed: boolean;

  constructor(
      startDate: string = '',
      completed: boolean = false,
  ) {
    this.startDate = startDate;
    this.completed = completed;
  }
}

let firstMealPlan = new MealPlan("01/01/1970");
let secondMealPlan = new MealPlan("01/08/1970", true);
let thirdMealPlan = new MealPlan("01/16/1970", true);
let fourthMealPlan = new MealPlan("02/02/1970", true);
let fifthMealPlan = new MealPlan("02/09/1970", true);

export default defineComponent({
  data() {
    return {
      mealPlans: [
        firstMealPlan,
        secondMealPlan,
        thirdMealPlan,
        fourthMealPlan,
        fifthMealPlan,
      ]
    }
  },
  methods: {
  },
});
</script>