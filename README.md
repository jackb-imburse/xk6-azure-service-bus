## Example xk6 build command with extension

xk6 build v0.39.0 --with github.com/jackb-imburse/xk6-azure-service-bus

## Example Usage
```
import asb from 'k6/x/azure-service-bus';

export default function () {
    var connectionString = '<CONNECTION_STRING>';
	var client = asb.getClient(connectionString);

    var queue = 'test-queue';
    var correlationId = 'TEST000000000:00000001';
    var subject = 'EXAMPLE_SUBJECT';
    var message = 'Test Message';

	asb.sendMessage(queue, message, client);
}
```