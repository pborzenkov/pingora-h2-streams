use clap::Parser;
use pingora_core::{apps::HttpServerOptions, server, upstreams::peer::HttpPeer, Result};
use pingora_proxy::Session;

#[derive(Parser, Debug, Default)]
pub struct ProxyOpt {
    #[clap(flatten)]
    server_opt: server::configuration::Opt,

    #[clap(short, long, help = "Listen address", default_value = "127.0.0.1:8081")]
    listen: String,

    #[clap(short, long, help = "Peer address", default_value = "127.0.0.1:8080")]
    peer: String,
    #[clap(long, help = "Max peer h2 streams", default_value = "100")]
    max_peer_streams: usize,
}

struct Proxy {
    peer: String,
    max_peer_streams: usize,
}

#[async_trait::async_trait]
impl pingora_proxy::ProxyHttp for Proxy {
    type CTX = ();

    fn new_ctx(&self) -> Self::CTX {}

    async fn upstream_peer(
        &self,
        _session: &mut Session,
        _ctx: &mut Self::CTX,
    ) -> Result<Box<HttpPeer>> {
        let mut peer = HttpPeer::new(&self.peer, false, "".into());
        peer.options.max_h2_streams = self.max_peer_streams;
        peer.options.set_http_version(2, 2);

        Ok(Box::new(peer))
    }
}

fn main() {
    let ProxyOpt {
        server_opt,
        listen,
        peer,
        max_peer_streams,
    } = ProxyOpt::parse();

    let mut server = server::Server::new(Some(server_opt)).unwrap();
    server.bootstrap();

    let mut proxy = pingora_proxy::http_proxy_service(
        &server.configuration,
        Proxy {
            peer,
            max_peer_streams,
        },
    );
    proxy.app_logic_mut().map(|p| {
        let mut opts = HttpServerOptions::default();
        opts.h2c = true;
        p.server_options = Some(opts);
    });
    proxy.add_tcp(&listen);
    server.add_service(proxy);

    server.run_forever();
}
