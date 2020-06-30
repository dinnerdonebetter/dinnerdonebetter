<template>
  <div class="registration-container">
    <div
      v-if="twoFactorQRCode === ''"
    >
      <el-form
        ref="registrationForm"
        :model="registrationForm"
        :rules="registrationRules"
        class="registration-form"
        autocomplete="on"
        label-position="left"
      >
        <el-form-item prop="username">
          <span class="svg-container">
            <svg-icon name="user" />
          </span>
          <el-input
            id="usernameInput"
            ref="username"
            v-model="registrationForm.username"
            name="username"
            type="text"
            autocomplete="on"
            placeholder="username"
          />
        </el-form-item>

        <el-form-item prop="password">
          <span class="svg-container">
            <svg-icon name="password" />
          </span>
          <el-input
            id="passwordInput"
            ref="password"
            v-model="registrationForm.password"
            name="password"
            type="password"
            placeholder="password"
            @keyup.enter.native="handleSubmit"
          />
        </el-form-item>

        <el-form-item prop="repeatedPassword">
          <span class="svg-container">
            <svg-icon name="password" />
          </span>
          <el-input
            id="passwordRepeatInput"
            ref="password"
            v-model="registrationForm.repeatedPassword"
            type="password"
            name="repeatedPassword"
            placeholder="password again"
            @keyup.enter.native="handleSubmit"
          />
        </el-form-item>

        <el-button
          id="registrationButton"
          :loading="loading"
          type="primary"
          style="width:100%; margin-bottom:30px;"
          @click.native.prevent="handleSubmit"
        >
          Register
        </el-button>
      </el-form>
    </div>
    <div
      v-else
      class="qr-container"
    >
      <img
        id="twoFactorSecretQRCode"
        style="width: 30%;"
        :src="twoFactorQRCode"
        alt="two factor authentication secret encoded as a QR code"
      >
      <p>
        You should save the secret this QR code contains, you'll be required to
        generate a token from it on every login.
      </p>
      <p>
        <input
          id="totpTokenInput"
          v-model="totpToken"
          type="text"
          placeholder="2FA Token"
          @keyup="shouldContinuationBeDisabled"
        >
      </p>
      <div v-if="twoFactorValidationError !== ''">
        {{ twoFactorValidationError }}
      </div>
      <p>
        <button
          id="totpTokenSubmitButton"
          :disabled="continueDisabled"
          @click="checkAndGoHome"
        >
          I've saved it!
        </button>
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import axios, { AxiosResponse } from 'axios';
import { Component, Vue } from 'vue-property-decorator';
import { Form as ElForm } from 'element-ui';
import { backendRoutes, statusCodes } from '@/constants';

@Component({
  name: 'Register',
})
export default class extends Vue {
  private readonly minimumPasswordCharacterCount = 8;
  private readonly requiredTOTPTokenLength = 6;

  private twoFactorValidationError = '';
  private twoFactorQRCode = '';
  private totpToken = '';
  private createdUserID: number | undefined = undefined;
  private loading = false;

  private continueDisabled = true;
  private shouldContinuationBeDisabled(): void {
    this.continueDisabled = isNaN(+this.totpToken) || this.totpToken.length !== this.requiredTOTPTokenLength;
  }

  private registrationForm = {
    username: '',
    password: '',
    repeatedPassword: '',
  }

  private registrationRules = {
    username: [{
      validator: (rule: object, value: string, callback: Function) => {
        if (value.length <= 0) {
          callback(new Error('Please enter a valid username'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    }],
    password: [{
      validator: (rule: object, value: string, callback: Function) => {
        if (value.length < this.minimumPasswordCharacterCount) {
          callback(new Error(`The password can not be less than ${this.minimumPasswordCharacterCount} digits`));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    }],
    repeatedPassword: [{
      validator: (rule: object, value: string, callback: Function) => {
        if (value === '') {
          callback(new Error('Please input the password again'));
        } else if (value !== this.registrationForm.password) {
          callback(new Error("passwords do not match!"));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    }],
  }

  private checkAndGoHome(): void {
    axios.post(backendRoutes.VERIFY_2FA_SECRET, {
        userID: this.createdUserID,
        totpToken: this.totpToken,
      }).then((response: AxiosResponse) => {
      if (response.status === statusCodes.ACCEPTED) {
        this.twoFactorValidationError = '';
        this.$router.push('/login');
      } else {
        console.log("bad response from server");
        this.twoFactorValidationError = `invalid response from 2FA validation route: ${response.status}`;
      }
    });
  }

  private handleSubmit(): void {
    (this.$refs.registrationForm as ElForm).validate(async(valid: boolean) => {
      if (valid) {
        axios.post(backendRoutes.USER_REGISTRATION, {
          username: this.registrationForm.username,
          password: this.registrationForm.password,
        }).then((response) => {
            if (response.status === statusCodes.CREATED) {
              return response.data;
            } else {
              throw new Error(`error registering new user: ${response.status}`);
            }
          })
          .then((data: {id: number; qrCode: string}) => {
            this.createdUserID = data["id"];
            this.twoFactorQRCode = data["qrCode"];
          });
      }
    });
  }
}
</script>

<style lang="scss">
// References: https://www.zhangxinxu.com/wordpress/2018/01/css-caret-color-first-line/
@supports (-webkit-mask: none) and (not (cater-color: $registrationCursorColor)) {
  .registration-container .el-input {
    input {
      color: $registrationCursorColor;
    }
    input::first-line {
      color: $lightGray;
    }
  }
}

.qr-container {
  display: block;
  text-align: center;
}

.qr-container img {
  display: block;
  margin: 0 auto;
}

.registration-container {
  .el-input {
    display: inline-block;
    height: 47px;
    width: 85%;

    input {
      height: 47px;
      background: transparent;
      border: 0px;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: $lightGray;
      caret-color: $registrationCursorColor;
      -webkit-appearance: none;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $registrationBackground inset !important;
        -webkit-text-fill-color: #fff !important;
      }
    }
  }

  .el-form-item {
    border: 1px solid rgba(255, 255, 255, 0.1);
    background: rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    color: #454545;
  }
}
</style>

<style lang="scss" scoped>
.registration-container {
  height: 100%;
  width: 100%;
  overflow: hidden;
  background-color: $registrationBackground;

  .registration-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $darkGray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $lightGray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $darkGray;
    cursor: pointer;
    user-select: none;
  }
}
</style>
