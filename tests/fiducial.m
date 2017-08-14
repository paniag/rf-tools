#!/usr/bin/octave -q
## fiducial.m
## Copyright 2017 Mac Radigan
## All Rights Reserved

   
  MESSAGE_SIZE           = 1024; % number of bytes in a single message
  SHORT_SIZE             = 256;  % number of residual bytes at end of stream
  NUMBER_OF_PATTERNS     = 15;   % number of messages in a sequence
  NUMBER_OF_REPETITIONS  = 5;    % number of times to repeat sequence (NB. forever if 0)

  args = argv();
  %% usage
  if(0 == size(args,1))
    fprintf(2, 'usage:  %s [MESSSAGE_SIZE] <SHORT_SIZE> <NUMBER_OF_PATTERNS> <NUMBER_OF_REPETITIONS>\n', mfilename);
    fprintf(2, '');
    fprintf(2, '  MESSSAGE_SIZE         - number of bytes sent per message [D:%d]\n',           MESSAGE_SIZE);
    fprintf(2, '  SHORT_SIZE            - number of residual bytes at end of stream [D:%d]\n',  SHORT_SIZE);
    fprintf(2, '  NUMBER_OF_PATTERNS    - number of distinct patterns to send [D:%d]\n',        NUMBER_OF_PATTERNS);
    fprintf(2, '  NUMBER_OF_REPETITIONS - number of times to repeat the sequence [D:%d]\n',     NUMBER_OF_REPETITIONS);
    exit(1);
  end

  %% parse command line arguments
  MESSAGE_SIZE = str2num(argv(){1});
  if(size(args,1) > 1)
    SHORT_SIZE = str2num(argv(){2});
  end
  if(size(args,1) > 2)
    NUMBER_OF_PATTERNS = str2num(argv(){3});
  end
  if(size(args,1) > 3)
    NUMBER_OF_REPETITIONS = str2num(argv(){4});
    if 0 == NUMBER_OF_REPETITIONS
      NUMBER_OF_REPETITIONS = Inf;
    end
  end

  %% generate message
  N_msg     = MESSAGE_SIZE;
  msg       = ones(1, N_msg);

  %% generate multiplier for each individual message
  N_pat     = NUMBER_OF_PATTERNS;
  pat       = 1:N_pat;

  %% create full sequence
  seq       = kron(msg, pat');

  %% print each repeated sequence
  while NUMBER_OF_REPETITIONS
    fprintf(1, dec2hex(seq')'); 
    NUMBER_OF_REPETITIONS = NUMBER_OF_REPETITIONS - 1;
  end

  %% print short sequence
  short_seq = repmat('X', 1, SHORT_SIZE);
  fprintf(1, short_seq);

## *EOF*
