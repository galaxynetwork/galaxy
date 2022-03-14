make install
make build
sudo rm -r ~/.galaxy
./bin/galaxyd init local
cp /Users/j/desktop/galaxies-labs/networks/galaxy-1/genesis.json ~/.galaxy/config
./bin/galaxyd add-genesis-account t1 1000000000000000uglx --keyring-backend os
./bin/galaxyd gentx t1 500000000000000uglx --chain-id galaxy-1 
./bin/galaxyd collect-gentxs
./bin/galaxyd start