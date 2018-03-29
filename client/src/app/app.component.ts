import {Component} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {CookieService} from 'ngx-cookie-service';
import {AuthService} from './auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';

  constructor(private http: HttpClient, private auth: AuthService) {
  }

  isLoggedIn(): boolean {
    return this.auth.isLoggedIn();
  }

  logout() {
    this.auth.logout();
  }

  onClickLogin() {
    this.auth.login();
  }

  onClickAuthrized() {
    this.http.get('/api/hello', {
      headers: this.auth.createAuthorizationHeader()
    }).subscribe(r => console.log(r));
  }

  onClickNotAuthrized() {
    this.http.get('/hello', {
      headers: this.auth.createAuthorizationHeader()
    }).subscribe(r => console.log(r));
  }
}
