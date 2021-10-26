<template>
  <div class="flex items-center justify-center h-screen px-6 bg-gray-200">
    <div class="w-full max-w-sm p-6 bg-white rounded-md shadow-md">
      <form class="mt-4" @submit.prevent="register" v-if="registrationVerificationQRCode === ''">
        <label class="block">
          <span class="text-sm text-gray-700">Username</span>
          <input
            type="text"
            class="block w-full mt-1 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
            placeholder="username"
            data-automation="registration_username_input"
            v-model="username"
          />
        </label>

        <label class="block mt-3">
          <span class="text-sm text-gray-700">Password</span>
          <input
            type="password"
            class="block w-full mt-1 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
            placeholder="something robust, please"
            data-automation="registration_password_input"
            v-model="password"
          />
        </label>

        <label class="block mt-3">
          <span class="text-sm text-gray-700">Password Again</span>
          <input
            type="password"
            class="block w-full mt-1 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
            placeholder="once more for the folks in the back"
            data-automation="registration_password_repeat_input"
            v-model="repeatedPassword"
          />
        </label>

        <div class="mt-6">
          <button
            type="submit"
            data-automation="register_button"
            class="w-full px-4 py-2 text-sm text-center text-white bg-indigo-600 rounded-md focus:outline-none hover:bg-indigo-500"
          >
            Register
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
          to="/login"
        >
          Login instead
        </router-link>
          </div>
        </div>
      </form>
      <form class="mt-4" @submit.prevent="confirmRegistration"  v-else>
        <img alt="qr encoded two factor secret" data-automation="two_factor_qr_code" class="w-full" :src="registrationVerificationQRCode" />

        <label class="block mt-3">
          <span class="text-sm text-gray-700">2FA Code</span>
          <input
              type="text"
              class="block w-full mt-1 border-gray-200 rounded-md focus:border-indigo-600 focus:ring focus:ring-opacity-40 focus:ring-indigo-500"
              placeholder="012345"
              data-automation="totp_secret_verification_token_input"
              v-model="verificationToken"
          />
        </label>

        <div class="mt-6">
          <button
              type="submit"
              data-automation="totp_token_submit_button"
              class="w-full px-4 py-2 text-sm text-center text-white bg-indigo-600 rounded-md focus:outline-none hover:bg-indigo-500"
          >
            Confirm Registration
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import axios, {AxiosError, AxiosResponse} from "axios";
import { defineComponent } from "vue";
import {settings} from "../../settings/settings";
import {backendRoutes} from "../../constants";

interface registrationResponse {
  qrCode: string;
  createdUserID: number;
}

export default defineComponent({
  data() {
    return {
      username: "",
      password: "",
      repeatedPassword: "",
      registrationVerificationQRCode: "",
      createdUserID: 0,
      verificationToken: "",
    }
  },
  methods: {
    register() {
      const registrationBody = {
        username: this.username,
        password: this.password,
      }

      axios.post(`${settings.API_SERVER_URL}${backendRoutes.USER_REGISTRATION}`, registrationBody)
        .then((result: AxiosResponse<registrationResponse>) => {
          this.registrationVerificationQRCode = result.data.qrCode;
          this.createdUserID = result.data.createdUserID;
        })
        .catch((err: AxiosError) => {
          console.error(err)
        });
    },
    confirmRegistration() {
      const confirmationBody = {
        totpToken: this.verificationToken,
        userID: this.createdUserID,
      }

      axios.post(`${settings.API_SERVER_URL}${backendRoutes.VERIFY_2FA_SECRET}`, confirmationBody)
          .then(() => {
            this.$router.push("/login");
          })
          .catch((err: AxiosError) => {
            console.error(err)
          });
    },
  },
});
</script>
