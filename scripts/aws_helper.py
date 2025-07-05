#!/usr/bin/env python3
"""
AWS CLI Helper Script
Provides a clean way to execute AWS CLI commands without shell interference
"""

import subprocess
import sys
import json
import os

def run_aws_command(args):
    """Run AWS CLI command with clean environment"""
    # Set clean environment
    env = os.environ.copy()
    env['PATH'] = '/usr/bin:/bin:/usr/sbin:/sbin:/opt/homebrew/bin'
    env['AWS_DEFAULT_REGION'] = 'af-south-1'
    
    # Build command
    cmd = ['/opt/homebrew/bin/aws'] + args
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, env=env)
        
        if result.returncode == 0:
            if result.stdout.strip():
                print(result.stdout.strip())
            return True
        else:
            print(f"Error: {result.stderr.strip()}")
            return False
            
    except Exception as e:
        print(f"Exception: {e}")
        return False

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 aws_helper.py <aws-command> [args...]")
        print("Example: python3 aws_helper.py sts get-caller-identity")
        sys.exit(1)
    
    # Get AWS command arguments
    aws_args = sys.argv[1:]
    
    # Run the command
    success = run_aws_command(aws_args)
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main() 