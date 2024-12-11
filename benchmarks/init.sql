CREATE EXTENSION plperlu;

CREATE OR REPLACE FUNCTION get_pid_cpu(int) returns float 
as
$$
 my $ps = "ps aux";
 my $pid = $_[0];
 my $awk = "awk '{if (\$2==" . $pid . "){print \$3}}'";
 my $cmd = $ps."|".$awk;
 $output = `$cmd`;
 my $cpu_perc = $output;
 return $cpu_perc;
$$ language plperlu;

CREATE OR REPLACE FUNCTION get_pid_mem(int) returns float 
as
$$
 my $ps = "ps aux";
 my $pid = $_[0];
 my $awk = "awk '{if (\$2==" . $pid . "){print \$4}}'";
 my $cmd = $ps."|".$awk;
 $output = `$cmd`;
 my $mem_perc = $output;
 return $mem_perc;
$$ language plperlu;
