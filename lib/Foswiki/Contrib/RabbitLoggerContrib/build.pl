#!/usr/bin/perl -w
use strict;
use warnings;

BEGIN { unshift @INC, split( /:/, $ENV{FOSWIKI_LIBS} ); }
use Foswiki::Contrib::Build;

package RabbitLoggerContribBuild;
our @ISA = qw(Foswiki::Contrib::Build);

sub new {
my $class = shift;
return bless($class->SUPER::new( "RabbitLoggerContrib" ), $class);
}

sub target_release {
my $this = shift;

print <<GUNK;

Building release $this->{RELEASE} of $this->{project}, from version $this->{VERSION}
GUNK
if ( $this->{-v} ) {
print 'Package name will be ', $this->{project}, "\n";
print 'Topic name will be ', $this->getTopicName(), "\n";
}

$this->_installDeps();

$this->build('compress');
$this->build('build');
$this->build('installer');
$this->build('stage');
$this->build('archive');
}

sub _installDeps {
    my $this = shift;

    local $| = 1;

    $this->pushd( "$this->{basedir}/dev/go/gologger" );

    print $this->sys_action( qw(go get) );

    my $toolsDir = "$this->{basedir}/tools";
    print $this->sys_action('mkdir', '-p', $toolsDir) unless -f $toolsDir;
    foreach my $cmd (qw(logtail logdump)) {
        $this->pushd( "$this->{basedir}/dev/go/gologger/cmd/$cmd" );
        print $this->sys_action('go', 'build', '-o', "$toolsDir/$cmd");
        $this->popd();
        if ($@ || !(-f "$toolsDir/$cmd")) {
            die "Could not build $cmd";
        }
    }

    my $resourcesDir = "$this->{basedir}/resources/gologger";
    foreach my $cmd (qw(logstore logreport)) {
        $this->pushd( "$this->{basedir}/dev/go/gologger/cmd/$cmd" );
        print $this->sys_action('go', 'build', '-o', "$resourcesDir/$cmd");
        $this->popd();
        if ($@ || !(-f "$resourcesDir/$cmd")) {
            die "Could not build $cmd";
        }
    }

    $this->popd();
}
my $build = RabbitLoggerContribBuild->new();
$build->build( $build->{target} );

