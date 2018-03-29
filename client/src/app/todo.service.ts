import {Injectable} from '@angular/core';
import {AuthService} from './auth.service';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {catchError} from 'rxjs/operators';
import {of} from 'rxjs/observable/of';
import {Todo} from './todo';

@Injectable()
export class TodoService {

  constructor(private http: HttpClient, private auth: AuthService) {
  }

  getList(): Observable<Todo[]> {
    return this.http.get('/api/todo/', {headers: this.auth.createAuthorizationHeader()})
      .pipe(
        catchError(this.handleError('getList', []))
      );
  }

  create(t: Todo): Observable<Todo> {
    return this.http.post('/api/todo/', t, {headers: this.auth.createAuthorizationHeader()})
      .pipe(
        catchError(this.handleError('create', []))
      );
  }

  update(t: Todo): Observable<Todo> {
    return this.http.put('/api/todo/', t, {headers: this.auth.createAuthorizationHeader()})
      .pipe(
        catchError(this.handleError('update', []))
      );
  }

  delete(t: Todo): Observable<Object> {
    return this.http.delete('/api/todo/' + t.id, {headers: this.auth.createAuthorizationHeader()})
      .pipe(
        catchError(this.handleError('delete', null))
      );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      console.error(error); // log to console instead
      return of(result as T);
    };
  }

}
