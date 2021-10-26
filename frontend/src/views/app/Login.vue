<template>
  <div class="flex items-center justify-center h-screen px-6 bg-gray-200">
    <div class="w-full max-w-sm p-6 bg-white rounded-md shadow-md">
      <form class="mt-4" @submit.prevent="login">
        <label class="block">
          <span class="text-sm text-gray-700">Username</span>
          <input
            type="text"
            class="block w-full mt-1 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
            placeholder="username"
            data-automation="login_username_input"
            v-model="username"
          />
        </label>

        <label class="block mt-3">
          <span class="text-sm text-gray-700">Password</span>
          <input
            type="password"
            class="block w-full mt-1 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
            placeholder="password"
            data-automation="login_password_input"
            v-model="password"
          />
        </label>

        <label class="block mt-3">
          <span class="text-sm text-gray-700">2FA Code</span>
          <input
            type="text"
            class="block w-full mt-1 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
            placeholder="012345"
            data-automation="login_totp_token_input"
            v-model="totpToken"
          />
        </label>

        <div class="mt-6">
          <button
            type="submit"
            data-automation="login_button"
            class="w-full px-4 py-2 text-sm text-center text-white bg-indigo-600 rounded-md focus:outline-none hover:bg-indigo-500"
          >
            Sign in
          </button>
        </div>

        <div class="flex items-center justify-between mt-4">
          <div>
            <label class="inline-flex items-center">
          <!--
              <input type="checkbox" class="text-indigo-600 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500" />
              <span class="mx-2 text-sm text-gray-600">Remember me</span>
          -->
            </label>
          </div>

          <div>

        <router-link
          class="block text-sm text-indigo-700 fontme hover:underline"
          to="/register"
        >
          Register instead
        </router-link>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import axios, {AxiosError, AxiosResponse} from "axios";
import { defineComponent } from "vue";

import { authenticationState, setAuthState } from "../../auth";
import { settings } from "../../settings/settings";
import {backendRoutes} from "../../constants";

export default defineComponent({
  data() {
    return {
      username: "",
      password: "",
      totpToken: "",
    }
  },
  methods: {
    login() {
      const loginBody = {
        "username": this.username,
        "password": this.password,
        "totpToken": this.totpToken,
      }

      axios.post(`${settings.API_SERVER_URL}${backendRoutes.LOGIN}`, loginBody)
          .then((result: AxiosResponse<authenticationState>) => {
            setAuthState(result.data);
            this.$router.push(result.data.userIsServiceAdmin ? "/admin/dashboard" : "/home");
          })
          .catch((err: AxiosError) => {
            console.error(err)
          });
    },
  },
});
</script>
