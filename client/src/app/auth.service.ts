import {Injectable} from '@angular/core';
import {CookieService} from 'ngx-cookie-service';
import {HttpHeaders} from '@angular/common/http';

@Injectable()
export class AuthService {

  static readonly tokenKey: string = 'TOKEN';

  token: string;

  constructor(private cookie: CookieService) {
    this.token = this.cookie.get(AuthService.tokenKey);
  }

  isLoggedIn(): boolean {
    return this.token !== '';
  }

  login() {
    location.href = '/oauth/start';
  }

  logout() {
    this.cookie.delete(AuthService.tokenKey);
    this.token = '';
  }

  createAuthorizationHeader(): HttpHeaders {
    return new HttpHeaders({'Authorization': 'Bearer ' + this.token});
  }
}
