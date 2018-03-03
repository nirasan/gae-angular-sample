import { Component } from '@angular/core';
import {HttpClient} from '@angular/common/http';

interface HelloResonse {
  Message: string
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';

  constructor(private http: HttpClient) {
    this.http.get<HelloResonse>('/hello?name=bob').subscribe(
      res => console.log(res)
    );
  }

  onClickLogin() {
    location.href = '/oauth/start';
  }
}
