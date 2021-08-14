{ lib, buildGoModule }:
buildGoModule rec {
  pname = "hochwassernachrichtendienst_exporter";
  version = "1.0.0";

  src = ./.;

  vendorSha256 = "sha256:0x81qssskbvby9ws26mpx028dd9zrmrzz6zxn732kmlg399rb8h0";
  runVend = true;

  meta = with lib; {
    description = "an unofficial prometheus exporter for the Hochwassernachrichtendienst Bayern";
    homepage = "https://github.com/mguentner/hochwassernachrichtendienst_exporter";
    license = licenses.mit;
    maintainers = with maintainers; [ mguentner ];
  };
}
