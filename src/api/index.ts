import { createClient } from '@connectrpc/connect';
import { GreetService } from './greet_pb';
import { createConnectTransport } from '@connectrpc/connect-web';
import { dev } from '$app/env';

const transport = createConnectTransport({
	// TODO: better way to handle this instead of hardcoding
	baseUrl: dev ? 'http://localhost:8080' : 'https://api-zoodoku.tizu.dev'
});

export const API = createClient(GreetService, transport);
export default API;
