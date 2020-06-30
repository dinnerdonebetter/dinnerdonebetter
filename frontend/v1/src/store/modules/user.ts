import axios, { AxiosResponse } from 'axios';
import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators';

import { getUserInfo } from '@/api/users';
import { backendRoutes, statusCodes } from '@/constants';
import { AuthStatus, LoginRequest } from '@/models/state';
import store from '@/store';

export interface UserState {
  isLoggedIn: boolean;
  isAdmin: boolean;
}

@Module({ dynamic: true, store, name: 'user', namespaced: true })
class User extends VuexModule implements UserState {
  public isLoggedIn = false;
  public isAdmin = false;

  @Mutation
  private setAuthStatus(authStatus: AuthStatus): void {
    this.isLoggedIn = authStatus.isAuthenticated;
    this.isAdmin = authStatus.isAdmin;
  }

  @Mutation
  private logOut(): void {
    this.isLoggedIn = false;
    this.isAdmin = false;
  }

  @Action({rawError: true})
  public async checkUserStatus() {
    return axios.get(backendRoutes.USER_AUTH_STATUS, {withCredentials: true})
      .then((statusResponse: AxiosResponse) => {
        this.context.commit("setAuthStatus", statusResponse.data as AuthStatus);
      })
      .catch((err) => {
        console.error(err);
      });
  }

  @Action({rawError: true})
  public async Login(creds: LoginRequest) {
    return axios.post(backendRoutes.LOGIN, creds, {withCredentials: true})
      .then(() => {
        getUserInfo().then((statusResponse: AxiosResponse<AuthStatus>) => {
          this.context.commit("setAuthStatus", statusResponse.data as AuthStatus);
        });
      });
  }

  @Action
  public async Logout() {
    if (!this.isLoggedIn) {
      throw Error('LogOut: already logged out!');
    }
    return axios.post(backendRoutes.LOGOUT, {
      withCredentials: true,
    }).then((response: AxiosResponse) => {
      if (response.status === statusCodes.OK) {
        this.context.commit("logOut");
      }
    });
  }
}

export const UserModule = getModule(User);
