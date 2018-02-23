import { Component } from '@angular/core';
import {HttpClient} from "@angular/common/http";

interface helloResonse {
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
    this.http.get<helloResonse>("/hello?name=bob").subscribe(
      res => console.log(res)
    );
  }
}
