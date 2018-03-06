import { Component } from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';
  token: string;

  constructor(private http: HttpClient, private cookie: CookieService) {
    this.token = this.cookie.get('TOKEN');
  }

  onClickLogin() {
    location.href = '/oauth/start';
  }

  onClickAuthrized() {
    this.http.get('/api/hello', {
      headers: new HttpHeaders({ 'Authorization': 'Bearer ' + this.token })
    }).subscribe(r => console.log(r));
  }

  onClickNotAuthrized() {
    this.http.get('/hello', {
      headers: new HttpHeaders({ 'Authorization': 'Bearer ' + this.token })
    }).subscribe(r => console.log(r));
  }
}
