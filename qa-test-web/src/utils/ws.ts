/***
 *  websocket interceptor
 * https://stackoverflow.com/questions/39210534/multiple-handlers-for-websocket-javascript
 * 
 */
import { websocket_url } from '../config'

export default function (path: string, protocol: string) {
  console.log(websocket_url)
  return new WebSocket(`${websocket_url}${path}`, protocol)
}