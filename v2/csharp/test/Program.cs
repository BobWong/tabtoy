using System;
using System.IO;
using tabtoy;
using table;
using Pis;

namespace test
{
    class MainClass
    {
        public static void Main(string[] args)
        {
            ConfigPool.Get<ConfigConfig>("/Users/bob/Workspace/go_workspace/src/github.com/davyxu/tabtoy/gen/GameConfig/Bytes");
        }
    }
}
