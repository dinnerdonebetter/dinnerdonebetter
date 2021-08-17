<template>
  <div class="grid gap-4 grid-cols-1 justify-items-center mt-3">
    <div class="w-1/2 bg-gray-50 p-6 rounded-lg shadow-lg">
      <h2 class="text-2xl font-bold mb-2 text-gray-800">Voting on your Meal Plan has finished!</h2>
      <p class="text-gray-700"><a href="#">Click here to see the results.</a></p>
    </div>

    <div class="w-1/2 flex flex-col bg-white shadow-lg rounded-lg overflow-hidden">
      <!-- <div class="bg-gray-200 text-gray-700 text-sm px-2 py-1">lasagnaluvr69 recommended a recipe</div> -->
      <div class="px-6 py-4">
        <h2 class="text-xl font-bold">I've done it again.</h2>
        <p>Check out my latest Lasagna opus magnus!</p>
        <img src="https://i.imgur.com/FpKhy1Q.jpeg">
      </div>

      <div class="bg-gray-200 px-6 py-4">
        <div class="flex items-center pt-3">
          <div class="bg-blue-700 w-12 h-12 flex justify-center items-center rounded-full">
          </div>
          <div class="ml-4">
            <p class="font-bold">lasangaluvr69</p>
            <p class="text-sm text-gray-700 mt-1">Lasagna Master</p>
          </div>
        </div>
      </div>
    </div>

    <div class="w-1/2 flex flex-col bg-white shadow-lg rounded-lg overflow-hidden">
      <div class="px-6 py-4 justify-center">
        <p>What is your favorite kind of rice?</p>
        <div>
          <input type="radio" id="arborio" name="drone" value="arborio">
          <label for="arborio">Arborio</label>
        </div>
        <div>
          <input type="radio" id="jasmine" name="drone" value="jasmine">
          <label for="jasmine">Jasmine</label>
        </div>
        <div>
          <input type="radio" id="basmati" name="drone" value="basmati">
          <label for="basmati">Basmati</label>
        </div>
        <div>
          <input type="radio" id="sushi" name="drone" value="sushi">
          <label for="sushi">Sushi</label>
        </div>
        <button class="py-3 px-6 text-white rounded-lg bg-purple-600 shadow-lg block md:inline-block uppercase">Submit</button>
      </div>

      <div class="bg-gray-200 px-6 py-4">
        <div class="flex items-center pt-3">
          <div class="bg-red-600 w-12 h-12 flex justify-center items-center rounded-full">
          </div>
          <div class="ml-4">
            <p class="font-bold">foodie_asshole</p>
            <p class="text-sm text-gray-700 mt-1">Tremendous Douche</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent} from "vue";

class Post {
  type: 'post' | 'meal_plan_vote_complete';
  title: string;

  constructor(
      title: string = "",
      type: 'post' | 'meal_plan_vote_complete' = 'post',
  ) {
    this.title = title;
    this.type = type;
  }
}

export default defineComponent({
  data() {
    return {
      posts: [
        new Post("made this cool salad!"),
        new Post("questions about bouef borgongongongong"),
        new Post("glenda posted a new recipe"),
      ]
    }
  },
  methods: {
  },
});
</script>