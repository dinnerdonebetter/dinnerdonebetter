<template>
  <div class="grid gap-4 grid-cols-1 justify-items-center mt-3">
    <div class="w-1/2 p-6 rounded-lg shadow-lg bg-gray-50" v-for="member in householdMembers">
      <img
          class="float-right w-8 h-8 overflow-hidden rounded-full shadow focus:outline-none object-cover w-full h-full"
          src="https://images.unsplash.com/photo-1528892952291-009c663ce843?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=296&q=80"
          alt="Your avatar"
      />
      <h2 class="text-2xl font-bold mb-2 text-gray-800">{{ member.name }}</h2>
    </div>
  </div>
</template>


<script lang="ts">
import { defineComponent} from "vue";

class HouseholdMember {
  name: string;

  constructor(
      name: string = '',
  ) {
    this.name = name;
  }
}

let firstHouseholdMember = new HouseholdMember("Jessie");
let secondHouseholdMember = new HouseholdMember("Becky");
let thirdHouseholdMember = new HouseholdMember("Joey");
let fourthHouseholdMember = new HouseholdMember("Jecky");
let fifthHouseholdMember = new HouseholdMember("Bowie");

export default defineComponent({
  data() {
    return {
      householdMembers: [
        firstHouseholdMember,
        secondHouseholdMember,
        thirdHouseholdMember,
        fourthHouseholdMember,
        fifthHouseholdMember,
      ]
    }
  },
  methods: {
  },
});
</script>