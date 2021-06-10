use echo::echoer_client::EchoerClient;
use echo::EchoRequest;

pub mod echo {
    tonic::include_proto!("echo"); // The string specified here must match the proto package name
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = EchoerClient::connect("http://127.0.0.1:5555").await?;

    let request = tonic::Request::new(EchoRequest {
        msg: "Hello Cole".into(),
    });

    let response = client.echo(request).await?;

    Ok(())
}
