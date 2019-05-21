# See bottom of file for license and copyright information
package Foswiki::Logger::RabbitMQLogger::EventIterator;
use strict;
use warnings;

our @ISA = qw/Foswiki::Iterator::EventIterator/;

package Foswiki::Logger::RabbitMQLogger;

use strict;
use warnings;

use JSON;

use Foswiki::Logger                           ();
use Foswiki::Iterator::EventIterator          ();
use Foswiki::Iterator::AggregateEventIterator ();
use Foswiki::Iterator::MergeEventIterator     ();
use Foswiki::Logger::RabbitMQLogger::Connection;
use Foswiki::Logger::RabbitMQLogger::StdErrDumper;
use Foswiki::Logger::RabbitMQLogger::Reader;
use Foswiki::Logger::RabbitMQLogger::EventIterator;

use Foswiki::Plugins::ModacHelpersPlugin;

our @ISA = ('Foswiki::Logger');

use Foswiki::Time qw(-nofoswiki);

use constant TRACE => 0;

sub new {
    my $class = shift;

    my $fallback = createSTDERRConnection();
    my $outbound;
    eval {
        $outbound = createRabbitMQConnection($fallback, _rabbitConfig());
    };
    if($@) {
        $outbound = $fallback;
        $outbound->send({fields => ['Error creating RabbitMQ connection', $@]});
    }
    my $this = {
        outbound => $outbound,
        fallback => $fallback,
    };
    return bless($this, $class);
}

sub log {
    my $this = shift;

    my $misc;
    my $level;
    my $extra;
    if(ref $_[0] eq 'HASH') {
        $misc = shift;
        $level = (delete $misc->{level}) || 'info';
        $extra = (delete $misc->{extra}) || [];
        $extra = [$extra] unless ref $extra eq 'ARRAY';
    } else {
        $level = shift;
        $extra = [@_];
        $misc = {};
    }

    my $skipStacktraceFrames = (delete $misc->{deleteStacktraceFrames}) || 0;
    my $noTrace = delete $misc->{noTrace};
    if((!$noTrace) && $level =~ m#^(?:warning|error|fatal)$#) {
        $misc->{tracestring} ||= _getTrace($skipStacktraceFrames + 1);
    }

    my $now  = (delete $misc->{time}) || time();
    my $time = Foswiki::Time::formatTime( $now, 'iso', 'servertime' );

    my $source = getSource();

    my $data = {
        metadata => {
        },
        log_data => {
            time => $time,
            level => $level,
            source => $source,
            fields => $extra,
            misc => $misc,
        },
    };

    if($level =~ m#^(?:warning|error|fatal)$# && (!$Foswiki::cfg{Extensions}{ModacHelpersPlugin}{NoBackendSentry}) && $extra->[0]) {
        $data->{metadata}->{report} = 'true';
        $data->{metadata}->{rms_user} = Foswiki::Plugins::ModacHelpersPlugin::getRmsCredentials();
        $data->{metadata}->{tags} = to_json({type => 'foswiki_backend'});
        $data->{metadata}->{environment} = $Foswiki::cfg{ModacHelpersPlugin}{Environment} || 'unknown_environment';
        $data->{metadata}->{release} = _getQwikiRelease();
    }

    eval {
        $this->{outbound}->send($data);
    };
    if($@) {
        use Data::Dumper;
        $this->{fallback}->send({fields => ["Error while sending log", $@, Dumper($data)]});
    }
}

sub getSource {
    return 'wiki:' . $Foswiki::cfg{DefaultUrlHost} =~ s#^https?://##r;
}

sub _getQwikiRelease {
    my $version = '(unknown)';
    eval {
        my $session = $Foswiki::Plugins::SESSION;
        if($session) {
            require Foswiki::Plugins::QueryVerionPlugin;
            $version = Foswiki::Plugins::QueryVerionPlugin::query($session, {_DEFAULT => 'QwikiContrib'});
        }
    };
    return "Q.wiki_$version";
}

sub _getTrace {
    my $skipFrames = shift || 0;

    my $traceString = '';
    eval {
        my $trace = Devel::StackTrace->new(
            skip_frames => $skipFrames + 2,
            max_arg_length => 50,
        );

        $traceString = $trace->as_string();
    };

    return $traceString;
}

sub eachEventSince {
    my ( $this, $time, $level, $version ) = @_;

    my ($startTime, $endTime);
    if(ref $time eq 'HASH') {
        $startTime = $time->{start};
        $endTime = $time->{end};
    } else {
        $startTime = $time;
    }

    my $source = getSource();
    my $events;
    eval {
        $events = Foswiki::Logger::RabbitMQLogger::Reader::read($source, $startTime, $endTime, $level, _rabbitConfig());
    };
    if($@) {
        $this->{fallback}->send({fields => ["Error while reading logs", $@]});
    }
    $events ||= [];

    return Foswiki::Logger::RabbitMQLogger::EventIterator->new($events, $version);
}

sub _rabbitConfig {
    my $rabbitHost = $Foswiki::cfg{Log}{RabbitMQLogger}{host} || 'localhost';
    my $rabbitUser = $Foswiki::cfg{Log}{RabbitMQLogger}{user} || 'guest';
    my $rabbitPassword = $Foswiki::cfg{Log}{RabbitMQLogger}{password} || 'guest';

    return ($rabbitHost, $rabbitUser, $rabbitPassword);
}

1;
__END__
Module of Foswiki - The Free and Open Source Wiki, http://foswiki.org/

Copyright (C) 2008-2013 Foswiki Contributors. Foswiki Contributors
are listed in the AUTHORS file in the root of this distribution.
NOTE: Please extend that file, not this notice.

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version. For
more details read LICENSE in the root of this distribution.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

As per the GPL, removal of this notice is prohibited.

